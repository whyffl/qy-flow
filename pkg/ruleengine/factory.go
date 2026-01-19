package ruleengine

import (
	"workflow/pkg/constant"
	"workflow/pkg/errors"
)

// NewRule 创建规则工厂方法
func NewRule(ruleType string, config RuleConfig, registry *RuleRegistry) (Rule, error) {
	switch ruleType {
	case constant.StringMatch:
		return &StringMatchRule{Config: config}, nil
	case constant.NumericCompare:
		return &NumericCompareRule{Config: config}, nil
	case constant.Handler:
		handlerName, ok := config.Value.(string)
		if !ok {
			return nil, errors.ErrInvalidFieldType(config.Field)
		}
		return &HandlerRule{Handler: handlerName, Registry: registry}, nil
	default:
		return nil, errors.ErrInvalidRuleType(ruleType)
	}
}

// NewRuleFromDefinition 从规则定义创建规则
func NewRuleFromDefinition(definition RuleDefinition, registry *RuleRegistry) (Rule, error) {
	return NewRule(definition.Type, definition.Config, registry)
}
