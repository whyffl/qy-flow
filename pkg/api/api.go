package api

import (
	"workflow/pkg/ruleengine"
	"workflow/pkg/workflowengine"
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
func (api *WorkflowEngineAPI) CreateWorkflow(workflow *workflowengine.Workflow) error {
	return api.engine.AddWorkflow(workflow)
}

// CreateWorkflowFromDefinition 从定义创建工作流
func (api *WorkflowEngineAPI) CreateWorkflowFromDefinition(def *workflowengine.WorkflowDefinition) error {
	return api.engine.AddWorkflowFromDefinition(def)
}

// ExecuteWorkflow 执行工作流
func (api *WorkflowEngineAPI) ExecuteWorkflow(flowKey string, data map[string]interface{}) (map[string]interface{}, error) {
	return api.engine.ExecuteWorkflow(flowKey, data)
}

// ExecuteWorkflowAsync 异步执行工作流
func (api *WorkflowEngineAPI) ExecuteWorkflowAsync(flowKey string, data map[string]interface{}, callback func(map[string]interface{}, error)) {
	api.engine.ExecuteWorkflowAsync(flowKey, data, callback)
}

// GetWorkflow 获取工作流
func (api *WorkflowEngineAPI) GetWorkflow(flowKey string) (*workflowengine.Workflow, error) {
	return api.engine.GetWorkflow(flowKey)
}

// DeleteWorkflow 删除工作流
func (api *WorkflowEngineAPI) DeleteWorkflow(flowKey string) error {
	return api.engine.RemoveWorkflow(flowKey)
}

// WorkflowEngineStats 工作流引擎统计信息
type WorkflowEngineStats struct {
	WorkflowCount int `json:"workflow_count"`
	RuleCount     int `json:"rule_count"`
}

// GetStats 获取统计信息
func (api *WorkflowEngineAPI) GetStats() *WorkflowEngineStats {
	return &WorkflowEngineStats{
		WorkflowCount: api.engine.GetWorkflowCount(),
		RuleCount:     api.engine.GetRuleEngine().GetRuleCount(),
	}
}

// RuleEngineAPI 规则引擎API
type RuleEngineAPI struct {
	engine *ruleengine.Engine
}

// NewRuleEngineAPI 创建规则引擎API
func NewRuleEngineAPI(engine *ruleengine.Engine) *RuleEngineAPI {
	return &RuleEngineAPI{
		engine: engine,
	}
}

// CreateRule 创建规则
func (api *RuleEngineAPI) CreateRule(ruleID string, ruleType string, config ruleengine.RuleConfig) error {
	rule, err := ruleengine.NewRule(ruleType, config)
	if err != nil {
		return err
	}

	// 如果是字符串匹配规则且使用正则，预编译
	if strRule, ok := rule.(*ruleengine.StringMatchRule); ok {
		if err := strRule.CompileRegex(); err != nil {
			return err
		}
	}

	return api.engine.AddRule(ruleID, rule)
}

// CreateRuleFromDefinition 从定义创建规则
func (api *RuleEngineAPI) CreateRuleFromDefinition(ruleID string, definition ruleengine.RuleDefinition) error {
	return api.engine.AddRuleFromDefinition(ruleID, definition)
}

// ExecuteRule 执行规则
func (api *RuleEngineAPI) ExecuteRule(ruleID string, data map[string]interface{}) (bool, error) {
	return api.engine.ExecuteRule(ruleID, data)
}

// ExecuteRulesConcurrent 并发执行多个规则
func (api *RuleEngineAPI) ExecuteRulesConcurrent(ruleIDs []string, data map[string]interface{}) (map[string]bool, error) {
	return api.engine.ExecuteRulesConcurrent(ruleIDs, data)
}

// ExecuteAllRules 执行所有规则
func (api *RuleEngineAPI) ExecuteAllRules(data map[string]interface{}) (map[string]bool, error) {
	return api.engine.ExecuteAllRules(data)
}

// DeleteRule 删除规则
func (api *RuleEngineAPI) DeleteRule(ruleID string) {
	api.engine.RemoveRule(ruleID)
}

// GetRule 获取规则
func (api *RuleEngineAPI) GetRule(ruleID string) (ruleengine.Rule, error) {
	return api.engine.GetRule(ruleID)
}

// GetRuleCount 获取规则数量
func (api *RuleEngineAPI) GetRuleCount() int {
	return api.engine.GetRuleCount()
}

// GetRegistry 获取规则注册表
func (api *RuleEngineAPI) GetRegistry() *ruleengine.RuleRegistry {
	return api.engine.GetRegistry()
}

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
