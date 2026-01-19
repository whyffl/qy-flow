package ruleengine

// Rule 规则接口
type Rule interface {
	// Execute 执行规则
	Execute(data map[string]interface{}) bool
	// GetType 获取规则类型
	GetType() string
}

// RuleConfig 规则配置
type RuleConfig struct {
	Field     string      `json:"field"`     // 字段名
	Operation string      `json:"operation"` // 操作符
	Value     interface{} `json:"value"`     // 值
}

// RuleDefinition 规则定义（用于JSON解析）
type RuleDefinition struct {
	RuleID string     `json:"rule_id"` // 规则ID
	Type   string     `json:"type"`    // 规则类型
	Config RuleConfig `json:"config"`  // 规则配置
}
