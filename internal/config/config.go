package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"workflow/pkg/ruleengine"
	"workflow/pkg/workflowengine"
)

// Config 应用配置
type Config struct {
	App            AppConfig            `yaml:"app"`
	WorkflowEngine WorkflowEngineConfig `yaml:"workflow_engine"`
	RuleEngine     RuleEngineConfig     `yaml:"rule_engine"`
}

// AppConfig 应用基础配置
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// WorkflowEngineConfig 工作流引擎配置（从 YAML 读取）
type WorkflowEngineConfig struct {
	MaxConcurrent   int    `yaml:"max_concurrent"`
	DefaultTimeout  int    `yaml:"default_timeout"`
	WorkerPoolSize  int    `yaml:"worker_pool_size"`
	EnableCache     bool   `yaml:"enable_cache"`
	CacheTTL        int    `yaml:"cache_ttl"`
	ShutdownTimeout string `yaml:"shutdown_timeout"`
}

// ToEngineConfig 转换为工作流引擎配置
func (c *WorkflowEngineConfig) ToEngineConfig() *workflowengine.Config {
	shutdownTimeout, err := time.ParseDuration(c.ShutdownTimeout)
	if err != nil {
		shutdownTimeout = 30 * time.Second
	}

	return &workflowengine.Config{
		MaxConcurrent:   c.MaxConcurrent,
		DefaultTimeout:  c.DefaultTimeout,
		WorkerPoolSize:  c.WorkerPoolSize,
		EnableCache:     c.EnableCache,
		CacheTTL:        c.CacheTTL,
		ShutdownTimeout: shutdownTimeout,
	}
}

// RuleEngineConfig 规则引擎配置（从 YAML 读取）
type RuleEngineConfig struct {
	DefaultTimeout int `yaml:"default_timeout"`
	WorkerPoolSize int `yaml:"worker_pool_size"`
	RegexPoolSize  int `yaml:"regex_pool_size"`
}

// ToEngineConfig 转换为规则引擎配置
func (c *RuleEngineConfig) ToEngineConfig() *ruleengine.Config {
	return &ruleengine.Config{
		DefaultTimeout: c.DefaultTimeout,
		WorkerPoolSize: c.WorkerPoolSize,
		RegexPoolSize:  c.RegexPoolSize,
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证并初始化配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &cfg, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("应用名称不能为空")
	}
	if c.App.Version == "" {
		return fmt.Errorf("应用版本不能为空")
	}

	if c.WorkflowEngine.MaxConcurrent <= 0 {
		return fmt.Errorf("工作流引擎最大并发数必须大于0")
	}
	if c.WorkflowEngine.DefaultTimeout <= 0 {
		return fmt.Errorf("工作流引擎默认超时时间必须大于0")
	}
	if c.WorkflowEngine.WorkerPoolSize <= 0 {
		return fmt.Errorf("工作流引擎工作池大小必须大于0")
	}
	if c.WorkflowEngine.CacheTTL <= 0 {
		return fmt.Errorf("工作流引擎缓存TTL必须大于0")
	}

	if c.RuleEngine.DefaultTimeout <= 0 {
		return fmt.Errorf("规则引擎默认超时时间必须大于0")
	}
	if c.RuleEngine.WorkerPoolSize <= 0 {
		return fmt.Errorf("规则引擎工作池大小必须大于0")
	}
	if c.RuleEngine.RegexPoolSize <= 0 {
		return fmt.Errorf("规则引擎正则池大小必须大于0")
	}

	return nil
}
