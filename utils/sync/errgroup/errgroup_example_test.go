package errgroup_test

import (
	"context"

	"github.com/xieziyu/go-coco/utils/sync/errgroup"
)

var job = func(ctx context.Context) error { return nil }

// 忽略任务失败, 所有任务完成才结束, 可以当做 sync.WaitGroup 使用
func ExampleGroup() {
	ctx := context.Background()
	var g errgroup.Group
	g.Go(func() error {
		return job(ctx)
	})
	g.Go(func() error {
		return job(ctx)
	})
	_ = g.Wait()
}

// 一旦一个任务失败, 自动 cancel 其他任务
func ExampleWithContext() {
	ctx := context.Background()
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return job(gctx)
	})
	g.Go(func() error {
		return job(gctx)
	})
	_ = g.Wait()
}
