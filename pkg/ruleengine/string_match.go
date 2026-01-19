package ruleengine

import (
	"regexp"
	"strings"

	"workflow/pkg/constant"
	"workflow/pkg/errors"
)

// StringMatchRule 字符串匹配规则
type StringMatchRule struct {
	Config RuleConfig
	regex  *regexp.Regexp
}

// Execute 执行字符串匹配规则
func (r *StringMatchRule) Execute(data map[string]interface{}) bool {
	fieldValue, exists := data[r.Config.Field]
	if !exists {
		return false
	}

	fieldStr, ok := fieldValue.(string)
	if !ok {
		return false
	}

	switch r.Config.Operation {
	case constant.Equals:
		return r.executeEquals(fieldStr)
	case constant.NotEquals:
		return r.executeNotEquals(fieldStr)
	case constant.Contains:
		return r.executeContains(fieldStr)
	case constant.NotContains:
		return r.executeNotContains(fieldStr)
	case constant.StartsWith:
		return r.executeStartsWith(fieldStr)
	case constant.EndsWith:
		return r.executeEndsWith(fieldStr)
	case constant.Regex:
		return r.executeRegex(fieldStr)
	default:
		return false
	}
}

// GetType 获取规则类型
func (r *StringMatchRule) GetType() string {
	return constant.StringMatch
}

// executeEquals 执行等于操作
func (r *StringMatchRule) executeEquals(fieldStr string) bool {
	valueStr, ok := r.Config.Value.(string)
	if !ok {
		return false
	}
	return fieldStr == valueStr
}

// executeNotEquals 执行不等于操作
func (r *StringMatchRule) executeNotEquals(fieldStr string) bool {
	valueStr, ok := r.Config.Value.(string)
	if !ok {
		return false
	}
	return fieldStr != valueStr
}

// executeContains 执行包含操作
func (r *StringMatchRule) executeContains(fieldStr string) bool {
	valueStr, ok := r.Config.Value.(string)
	if !ok {
		return false
	}
	return strings.Contains(fieldStr, valueStr)
}

// executeNotContains 执行不包含操作
func (r *StringMatchRule) executeNotContains(fieldStr string) bool {
	valueStr, ok := r.Config.Value.(string)
	if !ok {
		return false
	}
	return !strings.Contains(fieldStr, valueStr)
}

// executeStartsWith 执行以...开头操作
func (r *StringMatchRule) executeStartsWith(fieldStr string) bool {
	valueStr, ok := r.Config.Value.(string)
	if !ok {
		return false
	}
	return strings.HasPrefix(fieldStr, valueStr)
}

// executeEndsWith 执行以...结尾操作
func (r *StringMatchRule) executeEndsWith(fieldStr string) bool {
	valueStr, ok := r.Config.Value.(string)
	if !ok {
		return false
	}
	return strings.HasSuffix(fieldStr, valueStr)
}

// executeRegex 执行正则匹配操作
func (r *StringMatchRule) executeRegex(fieldStr string) bool {
	// 如果正则表达式未编译，先编译
	if r.regex == nil {
		pattern, ok := r.Config.Value.(string)
		if !ok {
			return false
		}
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return false
		}
		r.regex = compiled
	}
	return r.regex.MatchString(fieldStr)
}

// CompileRegex 预编译正则表达式（用于优化性能）
func (r *StringMatchRule) CompileRegex() error {
	if r.Config.Operation != constant.Regex {
		return nil
	}

	pattern, ok := r.Config.Value.(string)
	if !ok {
		return errors.ErrInvalidFieldType(r.Config.Field)
	}

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return errors.ErrRegexCompileError(pattern, err)
	}

	r.regex = compiled
	return nil
}
