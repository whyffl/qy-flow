package workflowengine

import (
	"time"
)

// Config 工作流引擎配置
type Config struct {
	MaxConcurrent   int           // 最大并发工作流数量
	DefaultTimeout  int           // 默认超时时间（秒）
	WorkerPoolSize  int           // 工作协程池大小
	EnableCache     bool          // 是否启用缓存
	CacheTTL        int           // 缓存过期时间（秒）
	ShutdownTimeout time.Duration // 关闭超时时间
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		MaxConcurrent:   1000,
		DefaultTimeout:  30,
		WorkerPoolSize:  100,
		EnableCache:     true,
		CacheTTL:        3600,
		ShutdownTimeout: 30 * time.Second,
	}
}

// Initialize 初始化配置，应用默认值
func (c *Config) Initialize() {
	if c.MaxConcurrent <= 0 {
		c.MaxConcurrent = 1000
	}
	if c.DefaultTimeout <= 0 {
		c.DefaultTimeout = 30
	}
	if c.WorkerPoolSize <= 0 {
		c.WorkerPoolSize = 100
	}
	if c.CacheTTL <= 0 {
		c.CacheTTL = 3600
	}
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = 30 * time.Second
	}
}

// GetTimeoutInSeconds 获取超时时间（秒）
func (c *Config) GetTimeoutInSeconds() int {
	return c.DefaultTimeout
}

// GetTimeoutDuration 获取超时时间（Duration）
func (c *Config) GetTimeoutDuration() time.Duration {
	return time.Duration(c.DefaultTimeout) * time.Second
}
