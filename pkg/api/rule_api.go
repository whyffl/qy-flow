package api

import (
	"workflow/pkg/ruleengine"
)

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
	rule, err := ruleengine.NewRule(ruleType, config, api.engine.GetRegistry())
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

// RegisterHandler 注册自定义处理器
func (api *RuleEngineAPI) RegisterHandler(name string, handler func(map[string]interface{}) bool) error {
	return api.engine.GetRegistry().RegisterHandler(name, handler)
}

// UnregisterHandler 注销自定义处理器
func (api *RuleEngineAPI) UnregisterHandler(name string) {
	api.engine.GetRegistry().UnregisterHandler(name)
}

// RuleEngineStats 规则引擎统计信息
type RuleEngineStats struct {
	RuleCount int `json:"rule_count"`
}

// GetStats 获取统计信息
func (api *RuleEngineAPI) GetStats() *RuleEngineStats {
	return &RuleEngineStats{
		RuleCount: api.engine.GetRuleCount(),
	}
}
