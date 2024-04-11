package app

// App 可部署的应用接口
type App interface {
	Start() error    // 应用的启动入口，包含初始化流程
	IsEnabled() bool // 应用启动的检测条件，例如 consumer 的环境变量开关是否开启等
}

// AlwaysEnabled 可用于嵌入，标识总是启用的应用
type AlwaysEnabled struct{}

func (a *AlwaysEnabled) IsEnabled() bool {
	return true
}
