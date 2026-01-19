package ruleengine

import "workflow/pkg/constant"

// NumericCompareRule 数值比较规则
type NumericCompareRule struct {
	Config RuleConfig
}

// Execute 执行数值比较规则
func (r *NumericCompareRule) Execute(data map[string]interface{}) bool {
	fieldValue, exists := data[r.Config.Field]
	if !exists {
		return false
	}

	fieldFloat, ok := toFloat64(fieldValue)
	if !ok {
		return false
	}

	value, ok := toFloat64(r.Config.Value)
	if !ok {
		return false
	}

	switch r.Config.Operation {
	case constant.Equal:
		return r.executeEqual(fieldFloat, value)
	case constant.NotEqual:
		return r.executeNotEqual(fieldFloat, value)
	case constant.GreaterThan:
		return r.executeGreaterThan(fieldFloat, value)
	case constant.GreaterOrEqual:
		return r.executeGreaterOrEqual(fieldFloat, value)
	case constant.LessThan:
		return r.executeLessThan(fieldFloat, value)
	case constant.LessOrEqual:
		return r.executeLessOrEqual(fieldFloat, value)
	default:
		return false
	}
}

// GetType 获取规则类型
func (r *NumericCompareRule) GetType() string {
	return constant.NumericCompare
}

// executeEqual 执行等于操作
func (r *NumericCompareRule) executeEqual(fieldFloat, value float64) bool {
	return fieldFloat == value
}

// executeNotEqual 执行不等于操作
func (r *NumericCompareRule) executeNotEqual(fieldFloat, value float64) bool {
	return fieldFloat != value
}

// executeGreaterThan 执行大于操作
func (r *NumericCompareRule) executeGreaterThan(fieldFloat, value float64) bool {
	return fieldFloat > value
}

// executeGreaterOrEqual 执行大于等于操作
func (r *NumericCompareRule) executeGreaterOrEqual(fieldFloat, value float64) bool {
	return fieldFloat >= value
}

// executeLessThan 执行小于操作
func (r *NumericCompareRule) executeLessThan(fieldFloat, value float64) bool {
	return fieldFloat < value
}

// executeLessOrEqual 执行小于等于操作
func (r *NumericCompareRule) executeLessOrEqual(fieldFloat, value float64) bool {
	return fieldFloat <= value
}

// toFloat64 将interface{}转换为float64
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case float64:
		return val, true
	case float32:
		return float64(val), true
	default:
		return 0, false
	}
}
