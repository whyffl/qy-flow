package executors

import (
	"context"
	"workflow/pkg/workflowengine/types"
)

// EndNodeExecutor 结束节点执行器
type EndNodeExecutor struct{}

// Execute 执行结束节点
func (e *EndNodeExecutor) Execute(ctx context.Context, node *types.Node, workflowCtx *types.WorkflowContext, workflow *types.Workflow) (*types.WorkflowContext, error) {
	_ = ctx
	_ = node
	_ = workflow

	// 设置下一个节点为0，表示工作流结束
	workflowCtx.SetNextNode(0)
	return workflowCtx, nil
}
