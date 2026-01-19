package executors

import (
	"context"
	"workflow/pkg/errors"
	"workflow/pkg/workflowengine/types"
	"workflow/pkg/workflowengine/utils"
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
func (e *SubprocessNodeExecutor) Execute(ctx context.Context, node *types.Node, workflowCtx *types.WorkflowContext, workflow *types.Workflow) (*types.WorkflowContext, error) {
	_ = workflow

	if node.SubprocessKey == "" {
		return nil, errors.ErrInvalidConfig("subprocess key is required")
	}

	// 检查子流程是否存在
	if _, err := e.engine.GetWorkflow(node.SubprocessKey); err != nil {
		return nil, err
	}

	// 创建子流程上下文
	subWorkflowCtx := &types.WorkflowContext{
		Data:     utils.DeepCopyMap(workflowCtx.GetData()),
		Metadata: make(map[string]interface{}),
	}

	// 执行子流程
	subResult, err := e.engine.ExecuteWorkflowWithContext(ctx, node.SubprocessKey, subWorkflowCtx)
	if err != nil {
		return nil, err
	}

	// 合并子流程结果到业务数据中
	for key, value := range subResult {
		workflowCtx.GetData()[key] = value
	}

	// 设置下一个节点（在上下文中，而不是修改业务数据）
	workflowCtx.SetNextNode(node.Next)
	return workflowCtx, nil
}
