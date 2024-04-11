package infra

// Infra 带有初始化逻辑的基础设施等
type Infra interface {
	CocoInfra() // 用于标记特定的接口实现方式
}

// infra.Coco 可用于嵌入，用于自定义的 infra 实例
type Coco struct{}

func (i Coco) CocoInfra() {}
