package types

import (
	"workflow/pkg/ruleengine"
)

// Case 节点中的条件分支
type Case struct {
	ConditionRuleID string `json:"condition_rule_id"` // 条件规则ID
	Next            int    `json:"next"`              // 下一个节点ID
}

// Node 工作流节点
type Node struct {
	ID            int    `json:"id"`             // 节点ID
	Type          string `json:"type"`           // 节点类型
	Name          string `json:"name"`           // 节点名称
	Next          int    `json:"next"`           // 下一个节点ID
	Cases         []Case `json:"cases"`          // 条件分支（switch节点使用）
	Default       *Case  `json:"default"`        // 默认分支（switch节点使用）
	Handler       string `json:"handler"`        // 处理器名称
	ParallelNodes []int  `json:"parallel_nodes"` // 并行节点列表（parallel节点使用）
	SubprocessKey string `json:"subprocess_key"` // 子流程key（subprocess节点使用）
	Timeout       int    `json:"timeout"`        // 超时时间（秒）
	Async         bool   `json:"async"`          // 是否异步执行
	Retry         int    `json:"retry"`          // 重试次数
	Description   string `json:"description"`    // 节点描述
}

// Metadata 工作流元数据
type Metadata struct {
	ID      int    `json:"id"`       // 工作流ID
	FlowKey string `json:"flow_key"` // 工作流唯一标识
	Name    string `json:"name"`     // 工作流名称
	Version string `json:"version"`  // 版本号
	Status  string `json:"status"`   // 状态：draft/published/disabled
}

// Workflow 工作流
type Workflow struct {
	Metadata Metadata                   `json:"metadata"` // 元数据
	Nodes    []Node                     `json:"nodes"`    // 节点列表
	Edges    []Edge                     `json:"edges"`    // 边列表（用于可视化）
	Rules    map[string]ruleengine.Rule `json:"-"`        // 规则映射（不序列化到JSON）
}

// Edge 边，用于连接节点（用于可视化）
type Edge struct {
	Source int `json:"source"` // 源节点ID
	Target int `json:"target"` // 目标节点ID
}

// WorkflowDefinition 工作流定义（用于JSON解析）
type WorkflowDefinition struct {
	Metadata Metadata            `json:"metadata"`
	Nodes    []Node              `json:"nodes"`
	Edges    []Edge              `json:"edges"`
	Rules    []RuleDefinitionRef `json:"rules"`
}

// RuleDefinitionRef 规则定义引用
type RuleDefinitionRef struct {
	RuleID string                `json:"rule_id"`
	Type   string                `json:"type"`
	Config ruleengine.RuleConfig `json:"config"`
}
