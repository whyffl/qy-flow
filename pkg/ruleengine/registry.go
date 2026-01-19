package ruleengine

import "fmt"

// RuleRegistry 规则注册表（用于管理自定义处理器）
type RuleRegistry struct {
	handlers map[string]func(map[string]interface{}) bool
}

// NewRuleRegistry 创建规则注册表
func NewRuleRegistry() *RuleRegistry {
	return &RuleRegistry{
		handlers: make(map[string]func(map[string]interface{}) bool),
	}
}

// RegisterHandler 注册自定义处理器
func (r *RuleRegistry) RegisterHandler(name string, handler func(map[string]interface{}) bool) error {
	if name == "" {
		return fmt.Errorf("handler name cannot be empty")
	}
	if handler == nil {
		return fmt.Errorf("handler function cannot be nil")
	}

	r.handlers[name] = handler
	return nil
}

// UnregisterHandler 注销自定义处理器
func (r *RuleRegistry) UnregisterHandler(name string) {
	delete(r.handlers, name)
}

// ExecuteHandler 执行自定义处理器
func (r *RuleRegistry) ExecuteHandler(name string, data map[string]interface{}) (bool, error) {
	handler, exists := r.handlers[name]
	if !exists {
		return false, fmt.Errorf("handler not found: %s", name)
	}

	return handler(data), nil
}

// HandlerExists 检查处理器是否存在
func (r *RuleRegistry) HandlerExists(name string) bool {
	_, exists := r.handlers[name]
	return exists
}
