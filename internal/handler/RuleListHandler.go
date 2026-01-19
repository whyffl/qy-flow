package handler

import "workflow/pkg/constant"

type RuleListHandlerRule struct {
	Handler string
}

// Execute 执行处理器规则
func (r *RuleListHandlerRule) Execute(data map[string]interface{}) bool {
	// 这里可以根据 handler 名称执行对应的 handler
	// 实际应用中，可以通过handler注册机制来执行自定义逻辑
	// 这里简单返回 true 表示执行成功
	// TODO: 实现自定义handler注册和执行机制
	return true
}

// GetType 获取规则类型
func (r *RuleListHandlerRule) GetType() string {
	return constant.Handler
}

// GetHandler 获取处理器名称
func (r *RuleListHandlerRule) GetHandler() string {
	return r.Handler
}
