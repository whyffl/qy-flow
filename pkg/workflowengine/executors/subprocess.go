package executors

import (
	"context"
	"workflow/pkg/errors"
	"workflow/pkg/workflowengine/types"
)

// SubprocessNodeExecutor 子流程节点执行器
type SubprocessNodeExecutor struct {
	engine types.EngineProvider
}

// NewSubprocessNodeExecutor 创建子流程节点执行器
func NewSubprocessNodeExecutor(engine types.EngineProvider) *SubprocessNodeExecutor {
	return &SubprocessNodeExecutor{
		engine: engine,
	}
}

// Execute 执行子流程节点
func (e *SubprocessNodeExecutor) Execute(ctx context.Context, node *types.Node, data map[string]interface{}, workflow *types.Workflow) (map[string]interface{}, error) {
	_ = workflow

	if node.SubprocessKey == "" {
		return nil, errors.ErrInvalidConfig("subprocess key is required")
	}

	// 检查子流程是否存在
	if _, err := e.engine.GetWorkflow(node.SubprocessKey); err != nil {
		return nil, err
	}

	// 执行子流程
	subResult, err := e.engine.ExecuteWorkflowWithContext(ctx, node.SubprocessKey, data)
	if err != nil {
		return nil, err
	}

	// 合并子流程结果
	for key, value := range subResult {
		data[key] = value
	}

	data["_next_node_id"] = node.Next
	return data, nil
}
