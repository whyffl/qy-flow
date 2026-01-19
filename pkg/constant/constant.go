package constant

// 规则类型常量
const (
	// StringMatch 字符串匹配规则类型
	StringMatch = "string_match"
	// NumericCompare 数值比较规则类型
	NumericCompare = "numeric_compare"
	// Handler 处理器规则类型
	Handler = "handler"
)

// handler 类型常量
const (
	RuleListHandler = "rule_list_handler"
)

// 字符串操作符常量
const (
	// Equals 等于
	Equals = "equals"
	// NotEquals 不等于
	NotEquals = "not_equals"
	// Contains 包含
	Contains = "contains"
	// NotContains 不包含
	NotContains = "not_contains"
	// StartsWith 以...开头
	StartsWith = "starts_with"
	// EndsWith 以...结尾
	EndsWith = "ends_with"
	// Regex 正则匹配
	Regex = "regex"
)

// 数值操作符常量
const (
	// Equal 数值等于
	Equal = "equal"
	// NotEqual 数值不等于
	NotEqual = "not_equal"
	// GreaterThan 大于
	GreaterThan = "greater_than"
	// GreaterOrEqual 大于等于
	GreaterOrEqual = "greater_or_equal"
	// LessThan 小于
	LessThan = "less_than"
	// LessOrEqual 小于等于
	LessOrEqual = "less_or_equal"
)

// 工作流节点类型常量
const (
	// Start 开始节点
	Start = "start"
	// Switch 分支节点
	Switch = "switch"
	// Action 动作节点
	Action = "action"
	// Parallel 并行节点
	Parallel = "parallel"
	// Subprocess 子流程节点
	Subprocess = "subprocess"
	// End 结束节点
	End = "end"
)

// 工作流状态常量
const (
	// StatusDraft 草稿状态
	StatusDraft = "draft"
	// StatusPublished 发布状态
	StatusPublished = "published"
	// StatusDisabled 禁用状态
	StatusDisabled = "disabled"
)
