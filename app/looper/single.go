// looper package provides utilities for looper job management.
package looper

import (
	"context"
	"syscall"
	"time"

	"github.com/iftechio/go-coco/utils/logger"
	"github.com/pkg/errors"
)

// 分布式锁
type DistributedLock interface {
	// Lock the key
	Lock(ctx context.Context, key string, expiration time.Duration) (ok bool, err error)
	// Unlock the key
	Unlock(ctx context.Context, key string) (err error)
	// Renew the key expiration time
	Renew(ctx context.Context, key string, expiration time.Duration) (err error)
}

// SinglePod 单Pod执行的任务
type SinglePod struct {
	ErrCh        chan error                      // 错误通道
	id           string                          // 唯一任务标识
	sleepTime    time.Duration                   // 每次任务间的执行间隔
	maxRunTime   time.Duration                   // 单次任务最大执行时间
	standbySleep time.Duration                   // 待机状态下休眠间隔
	run          func(ctx context.Context) error // 执行逻辑
	lock         DistributedLock                 // 分布式锁
	monitorTimer *time.Timer                     // 执行时长定时器: 处理 maxRunTime
	isTerminated bool                            // 是否已结束
}

// SinglePodConfig 任务配置
type SinglePodConfig struct {
	Sleep        time.Duration // 任务间的休眠时间
	MaxRunTime   time.Duration // [可选，默认300s] 单次任务最大执行时间
	StandbySleep time.Duration // [可选，默认10s] 待机状态下休眠时间
}

// NewSinglePod 创建单Pod执行的任务
func NewSinglePod(id string, lock DistributedLock, run func(ctx context.Context) error, cleanup func(), conf SinglePodConfig) (*SinglePod, func(), error) {
	if lock == nil {
		panic("NewSinglePod: lock is nil")
	}
	maxRunTime := conf.MaxRunTime
	if maxRunTime == 0 {
		// 默认 300s
		maxRunTime = 300 * time.Second
	}
	standbySleep := conf.StandbySleep
	if standbySleep == 0 {
		// 默认10s
		standbySleep = 10 * time.Second
	}
	job := &SinglePod{
		id:           id,
		lock:         lock,
		run:          run,
		ErrCh:        make(chan error, 1),
		sleepTime:    conf.Sleep,
		maxRunTime:   maxRunTime,
		standbySleep: standbySleep,
	}
	return job, func() {
		// 任务清除
		job.cleanup()
		// 外部依赖清除
		if cleanup != nil {
			cleanup()
		}
	}, nil
}

// getLockKey 获取锁key
func (j *SinglePod) getLockKey() string {
	return "single-pod-job:" + j.id
}

// Start 开始循环执行任务
func (j *SinglePod) Start() (err error) {
	for {
		if j.isTerminated {
			return
		}
		// 加锁
		ok, lerr := j.lock.Lock(context.Background(), j.getLockKey(), j.maxRunTime)
		if lerr != nil {
			return errors.WithStack(lerr)
		}
		if !ok {
			// standby mode
			logger.Debugf("job %s is standby", j.id)
			time.Sleep(j.standbySleep)
			continue
		}
		// working mode
		defer func() {
			_ = j.lock.Unlock(context.Background(), j.getLockKey())
		}()
		if err = j.working(); err != nil {
			return
		}
	}
}

// working 工作模式
func (j *SinglePod) working() (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("runtime error: %+v", err)
		}
	}()
	j.monitorTimer = time.NewTimer(j.maxRunTime)
	defer j.monitorTimer.Stop()
	for {
		if j.isTerminated {
			return
		}
		jobDone := make(chan bool)
		go func() {
			if !j.monitorTimer.Stop() {
				<-j.monitorTimer.C
			}
			j.monitorTimer.Reset(j.maxRunTime)
			select {
			case <-jobDone:
				return
			case <-j.monitorTimer.C:
				logger.Infof("job %s exceeded maxRunTime.", j.id)
				j.isTerminated = true
				_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
			}
		}()
		// 执行主逻辑
		if jerr := j.run(context.Background()); jerr != nil {
			j.ErrCh <- jerr
		}
		jobDone <- true
		// 延长锁
		if err = j.lock.Renew(context.Background(), j.getLockKey(), (j.maxRunTime + j.sleepTime*2)); err != nil {
			return errors.WithStack(err)
		}
		// 待机
		time.Sleep(j.sleepTime)
	}
}

// cleanup 清理
func (j *SinglePod) cleanup() {
	_ = j.lock.Unlock(context.Background(), j.getLockKey())
	if j.monitorTimer != nil {
		j.monitorTimer.Stop()
	}
	j.monitorTimer = nil
	logger.Infof("job %s cleanup", j.id)
}

// NotifyErr 运行时错误捕获
func (j *SinglePod) NotifyErr(fn func(err error)) {
	go func() {
		for {
			fn(<-j.ErrCh)
		}
	}()
}
