package executors

import (
	"context"
	"workflow/pkg/workflowengine/types"
)

// EndNodeExecutor 结束节点执行器
type EndNodeExecutor struct{}

// Execute 执行结束节点
func (e *EndNodeExecutor) Execute(ctx context.Context, node *types.Node, data map[string]interface{}, workflow *types.Workflow) (map[string]interface{}, error) {
	_ = ctx
	_ = node
	_ = workflow

	data["_next_node_id"] = 0
	return data, nil
}
