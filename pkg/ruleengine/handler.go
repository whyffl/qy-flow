package ruleengine

import "workflow/pkg/constant"

// HandlerRule 自定义处理器规则
type HandlerRule struct {
	Handler  string
	Registry *RuleRegistry
}

// Execute 执行处理器规则
func (r *HandlerRule) Execute(data map[string]interface{}) bool {
	if r.Registry == nil {
		return false
	}
	result, err := r.Registry.ExecuteHandler(r.Handler, data)
	if err != nil {
		return false
	}
	return result
}

// GetType 获取规则类型
func (r *HandlerRule) GetType() string {
	return constant.Handler
}

// GetHandler 获取处理器名称
func (r *HandlerRule) GetHandler() string {
	return r.Handler
}
