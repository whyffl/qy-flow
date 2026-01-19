package ruleengine

import (
	"context"
	"sync"
	"time"

	"workflow/pkg/errors"
)

// Engine 规则引擎
type Engine struct {
	rules      map[string]Rule
	ruleDefs   map[string]RuleDefinition
	registry   *RuleRegistry
	regexPool  *sync.Pool
	workerPool chan struct{}
	mu         sync.RWMutex
	timeout    time.Duration
	config     *Config
}

// NewEngine 创建规则引擎
func NewEngine(config *Config, registry *RuleRegistry) *Engine {
	if config == nil {
		config = DefaultConfig()
	}

	// 如果未传入 registry，则创建一个新的
	if registry == nil {
		registry = NewRuleRegistry()
	}

	engine := &Engine{
		rules:      make(map[string]Rule),
		ruleDefs:   make(map[string]RuleDefinition),
		registry:   registry,
		workerPool: make(chan struct{}, config.WorkerPoolSize),
		timeout:    config.GetTimeout(),
		config:     config,
	}

	// 初始化正则表达式池
	engine.regexPool = &sync.Pool{
		New: func() interface{} {
			return make([]*Rule, 0, config.RegexPoolSize)
		},
	}

	return engine
}

// AddRule 添加规则
func (e *Engine) AddRule(ruleID string, rule Rule) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if ruleID == "" {
		return errors.ErrInvalidConfig("rule ID cannot be empty")
	}

	e.rules[ruleID] = rule
	return nil
}

// AddRuleFromDefinition 从定义添加规则
func (e *Engine) AddRuleFromDefinition(ruleID string, definition RuleDefinition) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if ruleID == "" {
		return errors.ErrInvalidConfig("rule ID cannot be empty")
	}

	rule, err := NewRuleFromDefinition(definition, e.registry)
	if err != nil {
		return err
	}

	// 如果是字符串匹配规则且使用正则，预编译正则表达式
	if strRule, ok := rule.(*StringMatchRule); ok {
		if err := strRule.CompileRegex(); err != nil {
			return err
		}
	}

	e.rules[ruleID] = rule
	e.ruleDefs[ruleID] = definition
	return nil
}

// GetRule 获取规则
func (e *Engine) GetRule(ruleID string) (Rule, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	rule, exists := e.rules[ruleID]
	if !exists {
		return nil, errors.ErrFieldNotFound(ruleID)
	}

	return rule, nil
}

// RemoveRule 移除规则
func (e *Engine) RemoveRule(ruleID string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.rules, ruleID)
	delete(e.ruleDefs, ruleID)
}

// ExecuteRule 执行单个规则
func (e *Engine) ExecuteRule(ruleID string, data map[string]interface{}) (bool, error) {
	e.mu.RLock()
	rule, exists := e.rules[ruleID]
	e.mu.RUnlock()

	if !exists {
		return false, errors.ErrFieldNotFound(ruleID)
	}

	return rule.Execute(data), nil
}

// ExecuteRulesConcurrent 并发执行多个规则
func (e *Engine) ExecuteRulesConcurrent(ruleIDs []string, data map[string]interface{}) (map[string]bool, error) {
	results := make(map[string]bool)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, ruleID := range ruleIDs {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			// 获取工作协程
			e.workerPool <- struct{}{}
			defer func() { <-e.workerPool }()

			result, err := e.ExecuteRule(id, data)
			if err != nil {
				result = false
			}

			mu.Lock()
			results[id] = result
			mu.Unlock()
		}(ruleID)
	}

	wg.Wait()
	return results, nil
}

// ExecuteRulesWithTimeout 带超时的规则执行
func (e *Engine) ExecuteRulesWithTimeout(ruleIDs []string, data map[string]interface{}, timeout time.Duration) (map[string]bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	results := make(map[string]bool)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, ruleID := range ruleIDs {
		select {
		case <-ctx.Done():
			return results, errors.ErrExecutionTimeout("rule engine")
		default:
			wg.Add(1)

			go func(id string) {
				defer wg.Done()

				done := make(chan bool, 1)
				go func() {
					result, err := e.ExecuteRule(id, data)
					if err != nil {
						result = false
					}
					done <- result
				}()

				select {
				case result := <-done:
					mu.Lock()
					results[id] = result
					mu.Unlock()
				case <-ctx.Done():
					return
				}
			}(ruleID)
		}
	}

	wg.Wait()
	return results, nil
}

// ExecuteAllRules 执行所有规则
func (e *Engine) ExecuteAllRules(data map[string]interface{}) (map[string]bool, error) {
	e.mu.RLock()
	ruleIDs := make([]string, 0, len(e.rules))
	for ruleID := range e.rules {
		ruleIDs = append(ruleIDs, ruleID)
	}
	e.mu.RUnlock()

	return e.ExecuteRulesConcurrent(ruleIDs, data)
}

// GetRegistry 获取规则注册表
func (e *Engine) GetRegistry() *RuleRegistry {
	return e.registry
}

// GetRuleCount 获取规则数量
func (e *Engine) GetRuleCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.rules)
}

// Clear 清空所有规则
func (e *Engine) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.rules = make(map[string]Rule)
	e.ruleDefs = make(map[string]RuleDefinition)
}
