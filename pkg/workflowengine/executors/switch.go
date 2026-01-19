package executors

import (
	"context"
	"workflow/pkg/errors"
	"workflow/pkg/workflowengine/types"
)

// SwitchNodeExecutor 分支节点执行器
type SwitchNodeExecutor struct{}

// Execute 执行分支节点
func (e *SwitchNodeExecutor) Execute(ctx context.Context, node *types.Node, data map[string]interface{}, workflow *types.Workflow) (map[string]interface{}, error) {
	_ = ctx // 参数声明但未使用，符合接口定义要求
	// 遍历所有条件分支，找到第一个匹配的规则
	for _, caseItem := range node.Cases {
		rule, exists := workflow.Rules[caseItem.ConditionRuleID]
		if !exists {
			return nil, errors.ErrFieldNotFound(caseItem.ConditionRuleID)
		}

		if rule.Execute(data) {
			// 设置下一个节点
			data["_next_node_id"] = caseItem.Next
			return data, nil
		}
	}

	// 如果没有匹配的条件，使用默认分支
	if node.Default != nil {
		data["_next_node_id"] = node.Default.Next
		return data, nil
	}

	// 没有匹配的条件也没有默认分支，工作流结束
	data["_next_node_id"] = 0
	return data, nil
}
