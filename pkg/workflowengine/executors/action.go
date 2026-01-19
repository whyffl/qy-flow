package executors

import (
	"context"
	"fmt"
	"workflow/pkg/workflowengine/types"
)

// ActionNodeExecutor 动作节点执行器
type ActionNodeExecutor struct{}

// Execute 执行动作节点
func (e *ActionNodeExecutor) Execute(ctx context.Context, node *types.Node, workflowCtx *types.WorkflowContext, workflow *types.Workflow) (*types.WorkflowContext, error) {
	_ = ctx
	_ = workflow

	if node.Handler == "" {
		workflowCtx.SetNextNode(node.Next)
		return workflowCtx, nil
	}

	// TODO: 实现handler注册和执行机制
	// 这里可以扩展为调用外部handler
	fmt.Printf("Executing action handler: %s\n", node.Handler)

	// 设置下一个节点（在上下文中，而不是修改业务数据）
	workflowCtx.SetNextNode(node.Next)
	return workflowCtx, nil
}
