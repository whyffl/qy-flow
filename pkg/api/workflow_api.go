package api

import (
	"context"

	"workflow/pkg/workflowengine"
	"workflow/pkg/workflowengine/types"
)

// WorkflowEngineAPI 工作流引擎API
type WorkflowEngineAPI struct {
	engine *workflowengine.Engine
}

// NewWorkflowEngineAPI 创建工作流引擎API
func NewWorkflowEngineAPI(engine *workflowengine.Engine) *WorkflowEngineAPI {
	return &WorkflowEngineAPI{
		engine: engine,
	}
}

// CreateWorkflow 创建工作流
func (api *WorkflowEngineAPI) CreateWorkflow(workflow *types.Workflow) error {
	return api.engine.AddWorkflow(workflow)
}

// CreateWorkflowFromDefinition 从定义创建工作流
func (api *WorkflowEngineAPI) CreateWorkflowFromDefinition(def *types.WorkflowDefinition) error {
	return api.engine.AddWorkflowFromDefinition(def)
}

// ExecuteWorkflowWithContext 带上下文执行工作流
func (api *WorkflowEngineAPI) ExecuteWorkflowWithContext(ctx context.Context, flowKey string, workflowCtx *types.WorkflowContext) (map[string]interface{}, error) {
	return api.engine.ExecuteWorkflowWithContext(ctx, flowKey, workflowCtx)
}

// ExecuteWorkflowAsync 异步执行工作流
func (api *WorkflowEngineAPI) ExecuteWorkflowAsync(ctx context.Context, flowKey string, workflowCtx *types.WorkflowContext, callback func(map[string]interface{}, error)) {
	api.engine.ExecuteWorkflowAsync(ctx, flowKey, workflowCtx, callback)
}

// GetWorkflow 获取工作流
func (api *WorkflowEngineAPI) GetWorkflow(flowKey string) (*types.Workflow, error) {
	return api.engine.GetWorkflow(flowKey)
}

// DeleteWorkflow 删除工作流
func (api *WorkflowEngineAPI) DeleteWorkflow(flowKey string) error {
	return api.engine.RemoveWorkflow(flowKey)
}

// WorkflowEngineStats 工作流引擎统计信息
type WorkflowEngineStats struct {
	WorkflowCount int `json:"workflow_count"`
}

// GetStats 获取统计信息
func (api *WorkflowEngineAPI) GetStats() *WorkflowEngineStats {
	return &WorkflowEngineStats{
		WorkflowCount: api.engine.GetWorkflowCount(),
	}
}
