package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/iftechio/go-coco/infra"
	"github.com/iftechio/go-coco/utils/logger"
	"github.com/iftechio/go-coco/utils/sync/errgroup"
)

type Manager struct {
	apps   []App
	infras []infra.Infra
}

func NewManager() Manager {
	return Manager{}
}

// RegisterApp 注册一个 App 实例
func (m *Manager) RegisterApp(app App) {
	if m.apps == nil {
		m.apps = make([]App, 0)
	}
	m.apps = append(m.apps, app)
}

// Run 执行所有Enabled App
func (m *Manager) Run(cleanup func()) (err error) {
	// 开始执行 App
	errCh := make(chan error, 1)
	go func() {
		errCh <- m.runApps()
	}()

	// make a channel to receive OS signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// wait until one of channel received
	select {
	case <-sigCh:
		// do nothing
	case err = <-errCh:
		if err != nil {
			logger.Error(err)
		}
	}

	if cleanup != nil {
		cleanup()
	}
	return
}

// runApps 执行所有 Enabled App
func (m *Manager) runApps() error {
	if len(m.apps) == 0 {
		logger.Info("no app to run")
		return nil
	}
	eg := errgroup.New()
	for _, app := range m.apps {
		if app.IsEnabled() {
			a := app
			eg.Go(a.Start)
		}
	}
	return eg.Wait()
}

// RegisterInfra 注册一个 Infra 实例
func (m *Manager) RegisterInfra(i infra.Infra) {
	if m.infras == nil {
		m.infras = make([]infra.Infra, 0)
	}
	m.infras = append(m.infras, i)
}
