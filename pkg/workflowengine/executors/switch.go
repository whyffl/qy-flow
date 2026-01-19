package executors

import (
	"context"
	"workflow/pkg/errors"
	"workflow/pkg/workflowengine/types"
)

// SwitchNodeExecutor 分支节点执行器
type SwitchNodeExecutor struct{}

// Execute 执行分支节点
func (e *SwitchNodeExecutor) Execute(ctx context.Context, node *types.Node, workflowCtx *types.WorkflowContext, workflow *types.Workflow) (*types.WorkflowContext, error) {
	_ = ctx // 参数声明但未使用，符合接口定义要求
	// 遍历所有条件分支，找到第一个匹配的规则
	for _, caseItem := range node.Cases {
		rule, exists := workflow.Rules[caseItem.ConditionRuleID]
		if !exists {
			return nil, errors.ErrFieldNotFound(caseItem.ConditionRuleID)
		}

		if rule.Execute(workflowCtx.GetData()) {
			// 设置下一个节点（在上下文中，而不是修改业务数据）
			workflowCtx.SetNextNode(caseItem.Next)
			return workflowCtx, nil
		}
	}

	// 如果没有匹配的条件，使用默认分支
	if node.Default != nil {
		workflowCtx.SetNextNode(node.Default.Next)
		return workflowCtx, nil
	}

	// 没有匹配的条件也没有默认分支，工作流结束
	workflowCtx.SetNextNode(0)
	return workflowCtx, nil
}
