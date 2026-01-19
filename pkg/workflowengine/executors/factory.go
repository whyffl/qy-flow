package executors

import (
	"workflow/pkg/constant"
	"workflow/pkg/workflowengine/types"
)

// nodeExecutorFactory 执行器工厂
type nodeExecutorFactory struct {
	engine types.EngineProvider
}

// newNodeExecutorFactory 创建执行器工厂
func newNodeExecutorFactory(engine types.EngineProvider) *nodeExecutorFactory {
	return &nodeExecutorFactory{
		engine: engine,
	}
}

// GetNodeExecutor 根据节点类型获取对应的执行器
func GetNodeExecutor(nodeType string) types.NodeExecutor {
	factory := &nodeExecutorFactory{}
	return factory.CreateExecutor(nodeType)
}

// GetNodeExecutorWithEngine 根据节点类型和引擎实例获取对应的执行器
func GetNodeExecutorWithEngine(nodeType string, engine types.EngineProvider) types.NodeExecutor {
	factory := newNodeExecutorFactory(engine)
	return factory.CreateExecutor(nodeType)
}

// CreateExecutor 创建执行器
func (f *nodeExecutorFactory) CreateExecutor(nodeType string) types.NodeExecutor {
	switch nodeType {
	case constant.Start:
		return &StartNodeExecutor{}
	case constant.Switch:
		return &SwitchNodeExecutor{}
	case constant.Action:
		return &ActionNodeExecutor{}
	case constant.Parallel:
		return &ParallelNodeExecutor{}
	case constant.Subprocess:
		return NewSubprocessNodeExecutor(f.engine)
	case constant.End:
		return &EndNodeExecutor{}
	default:
		// 默认使用 end 执行器
		return &EndNodeExecutor{}
	}
}
