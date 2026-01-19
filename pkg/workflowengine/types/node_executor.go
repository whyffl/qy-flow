package types

import "context"

// NodeExecutor 节点执行器接口
type NodeExecutor interface {
	Execute(ctx context.Context, node *Node, data map[string]interface{}, workflow *Workflow) (map[string]interface{}, error)
}
