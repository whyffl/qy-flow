package api

import (
	"workflow/pkg/ruleengine"
	"workflow/pkg/workflowengine"
)

// EngineAPI 组合API
type EngineAPI struct {
	Workflow *WorkflowEngineAPI
	Rule     *RuleEngineAPI
}

// NewEngineAPI 创建组合API
func NewEngineAPI(wfEngine *workflowengine.Engine, ruleEngine *ruleengine.Engine) *EngineAPI {
	return &EngineAPI{
		Workflow: NewWorkflowEngineAPI(wfEngine),
		Rule:     NewRuleEngineAPI(ruleEngine),
	}
}
