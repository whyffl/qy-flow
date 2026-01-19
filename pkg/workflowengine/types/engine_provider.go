package types

import "context"

// EngineProvider 引擎提供者接口（用于解决循环引用）
type EngineProvider interface {
	GetWorkflow(flowKey string) (*Workflow, error)
	ExecuteWorkflowWithContext(ctx context.Context, flowKey string, data map[string]interface{}) (map[string]interface{}, error)
}
