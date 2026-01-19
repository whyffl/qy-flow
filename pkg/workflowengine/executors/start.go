package executors

import (
	"context"
	"workflow/pkg/workflowengine/types"
)

// StartNodeExecutor 开始节点执行器
type StartNodeExecutor struct{}

// Execute 执行开始节点
func (e *StartNodeExecutor) Execute(ctx context.Context, node *types.Node, workflowCtx *types.WorkflowContext, workflow *types.Workflow) (*types.WorkflowContext, error) {
	// 参数声明但未使用，符合接口定义要求
	_ = ctx
	_ = workflow

	// 开始节点通常不做任何处理，只是传递数据并设置下一个节点
	workflowCtx.SetNextNode(node.Next)
	return workflowCtx, nil
}
