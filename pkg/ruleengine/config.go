package ruleengine

import (
	"time"
)

// Config 规则引擎配置
type Config struct {
	DefaultTimeout int // 默认超时时间（秒）
	WorkerPoolSize int // 工作协程池大小
	RegexPoolSize  int // 正则表达式池大小
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DefaultTimeout: 30,
		WorkerPoolSize: 50,
		RegexPoolSize:  100,
	}
}

// Initialize 初始化配置，应用默认值
func (c *Config) Initialize() {
	if c.DefaultTimeout <= 0 {
		c.DefaultTimeout = 30
	}
	if c.WorkerPoolSize <= 0 {
		c.WorkerPoolSize = 50
	}
	if c.RegexPoolSize <= 0 {
		c.RegexPoolSize = 100
	}
}

// GetTimeout 获取超时时间
func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.DefaultTimeout) * time.Second
}
