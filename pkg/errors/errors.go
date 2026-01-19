package errors

import "fmt"

// WorkflowError 工作流引擎错误
type WorkflowError struct {
	Code    string
	Message string
	Err     error
}

// Error 实现error接口
func (e *WorkflowError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 解包内部错误
func (e *WorkflowError) Unwrap() error {
	return e.Err
}

// NewWorkflowError 创建工作流错误
func NewWorkflowError(code, message string, err error) *WorkflowError {
	return &WorkflowError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// RuleEngineError 规则引擎错误
type RuleEngineError struct {
	Code    string
	Message string
	Err     error
}

// Error 实现error接口
func (e *RuleEngineError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 解包内部错误
func (e *RuleEngineError) Unwrap() error {
	return e.Err
}

// NewRuleEngineError 创建规则引擎错误
func NewRuleEngineError(code, message string, err error) *RuleEngineError {
	return &RuleEngineError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 预定义错误代码
const (
	// 错误代码：工作流
	ErrCodeWorkflowNotFound   = "WF001"
	ErrCodeNodeNotFound       = "WF002"
	ErrCodeInvalidNodeType    = "WF003"
	ErrCodeInvalidConfig      = "WF004"
	ErrCodeExecutionTimeout   = "WF005"
	ErrCodeExecutionFailed    = "WF006"
	ErrCodeParallelFailed     = "WF007"
	ErrCodeSubprocessNotFound = "WF008"
	ErrCodeInvalidJSON        = "WF009"

	// 错误代码：规则引擎
	ErrCodeInvalidRuleType   = "RE001"
	ErrCodeInvalidOperation  = "RE002"
	ErrCodeFieldNotFound     = "RE003"
	ErrCodeInvalidFieldType  = "RE004"
	ErrCodeRegexCompileError = "RE005"
)

// 创建预定义错误的辅助函数

// ErrWorkflowNotFound 创建工作流未找到错误
func ErrWorkflowNotFound(flowKey string) *WorkflowError {
	return NewWorkflowError(ErrCodeWorkflowNotFound, fmt.Sprintf("workflow not found: %s", flowKey), nil)
}

// ErrNodeNotFound 创建节点未找到错误
func ErrNodeNotFound(nodeID int) *WorkflowError {
	return NewWorkflowError(ErrCodeNodeNotFound, fmt.Sprintf("node not found: %d", nodeID), nil)
}

// ErrInvalidNodeType 创建无效节点类型错误
func ErrInvalidNodeType(nodeType string) *WorkflowError {
	return NewWorkflowError(ErrCodeInvalidNodeType, fmt.Sprintf("invalid node type: %s", nodeType), nil)
}

// ErrInvalidConfig 创建无效配置错误
func ErrInvalidConfig(detail string) *WorkflowError {
	return NewWorkflowError(ErrCodeInvalidConfig, fmt.Sprintf("invalid configuration: %s", detail), nil)
}

// ErrExecutionTimeout 创建执行超时错误
func ErrExecutionTimeout(flowKey string) *WorkflowError {
	return NewWorkflowError(ErrCodeExecutionTimeout, fmt.Sprintf("execution timeout: %s", flowKey), nil)
}

// ErrExecutionFailed 创建执行失败错误
func ErrExecutionFailed(nodeID int, err error) *WorkflowError {
	return NewWorkflowError(ErrCodeExecutionFailed, fmt.Sprintf("execution failed at node: %d", nodeID), err)
}

// ErrParallelFailed 创建并行执行失败错误
func ErrParallelFailed(err error) *WorkflowError {
	return NewWorkflowError(ErrCodeParallelFailed, "parallel execution failed", err)
}

// ErrSubprocessNotFound 创建子流程未找到错误
func ErrSubprocessNotFound(subprocessKey string) *WorkflowError {
	return NewWorkflowError(ErrCodeSubprocessNotFound, fmt.Sprintf("subprocess not found: %s", subprocessKey), nil)
}

// ErrInvalidJSON 创建JSON解析错误
func ErrInvalidJSON(err error) *WorkflowError {
	return NewWorkflowError(ErrCodeInvalidJSON, "invalid JSON format", err)
}

// ErrInvalidRuleType 创建无效规则类型错误
func ErrInvalidRuleType(ruleType string) *RuleEngineError {
	return NewRuleEngineError(ErrCodeInvalidRuleType, fmt.Sprintf("invalid rule type: %s", ruleType), nil)
}

// ErrInvalidOperation 创建无效操作符错误
func ErrInvalidOperation(operation string) *RuleEngineError {
	return NewRuleEngineError(ErrCodeInvalidOperation, fmt.Sprintf("invalid operation: %s", operation), nil)
}

// ErrFieldNotFound 创建字段未找到错误
func ErrFieldNotFound(fieldName string) *RuleEngineError {
	return NewRuleEngineError(ErrCodeFieldNotFound, fmt.Sprintf("field not found: %s", fieldName), nil)
}

// ErrInvalidFieldType 创建无效字段类型错误
func ErrInvalidFieldType(fieldName string) *RuleEngineError {
	return NewRuleEngineError(ErrCodeInvalidFieldType, fmt.Sprintf("invalid field type: %s", fieldName), nil)
}

// ErrRegexCompileError 创建正则表达式编译错误
func ErrRegexCompileError(pattern string, err error) *RuleEngineError {
	return NewRuleEngineError(ErrCodeRegexCompileError, fmt.Sprintf("regex compile error: %s", pattern), err)
}
