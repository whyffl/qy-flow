package workflowengine

import (
	"context"
	"sync"
	"time"
	"workflow/pkg/workflowengine/executors"

	"workflow/pkg/errors"
	"workflow/pkg/ruleengine"
	"workflow/pkg/workflowengine/types"
)

// Engine 工作流引擎
type Engine struct {
	workflows    map[string]*types.Workflow
	config       *Config
	ruleEngine   *ruleengine.Engine
	workerPool   chan struct{}
	semaphore    chan struct{}
	cache        *WorkflowCache
	mu           sync.RWMutex
	shutdownOnce sync.Once
}

// WorkflowCache 工作流缓存
type WorkflowCache struct {
	cache map[string]*cacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
}

type cacheEntry struct {
	workflow *types.Workflow
	expires  time.Time
}

// NewEngine 创建工作流引擎
func NewEngine(config *Config, ruleEngine *ruleengine.Engine) *Engine {
	if config == nil {
		config = DefaultConfig()
	}

	engine := &Engine{
		workflows:  make(map[string]*types.Workflow),
		config:     config,
		ruleEngine: ruleEngine,
		workerPool: make(chan struct{}, config.WorkerPoolSize),
		semaphore:  make(chan struct{}, config.MaxConcurrent),
	}

	// 初始化缓存
	if config.EnableCache {
		engine.cache = &WorkflowCache{
			cache: make(map[string]*cacheEntry),
			ttl:   time.Duration(config.CacheTTL) * time.Second,
		}
	}

	return engine
}

// AddWorkflow 添加工作流
func (e *Engine) AddWorkflow(workflow *types.Workflow) error {
	if workflow == nil {
		return errors.ErrInvalidConfig("workflow cannot be nil")
	}

	if workflow.Metadata.FlowKey == "" {
		return errors.ErrInvalidConfig("flow key cannot be empty")
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.workflows[workflow.Metadata.FlowKey] = workflow
	return nil
}

// AddWorkflowFromDefinition 从定义添加工作流
func (e *Engine) AddWorkflowFromDefinition(def *types.WorkflowDefinition) error {
	workflow := &types.Workflow{
		Metadata: def.Metadata,
		Nodes:    def.Nodes,
		Edges:    def.Edges,
		Rules:    make(map[string]ruleengine.Rule),
	}

	// 创建规则
	for _, ruleRef := range def.Rules {
		rule, err := ruleengine.NewRule(
			ruleRef.Type,
			ruleRef.Config,
			e.ruleEngine.GetRegistry(),
		)
		if err != nil {
			return err
		}

		// 如果是字符串匹配规则且使用正则，预编译
		if strRule, ok := rule.(*ruleengine.StringMatchRule); ok {
			if err := strRule.CompileRegex(); err != nil {
				return err
			}
		}

		workflow.Rules[ruleRef.RuleID] = rule
	}

	return e.AddWorkflow(workflow)
}

// GetWorkflow 获取工作流
func (e *Engine) GetWorkflow(flowKey string) (*types.Workflow, error) {
	// 先从缓存获取
	if e.cache != nil {
		if cached := e.cache.Get(flowKey); cached != nil {
			return cached, nil
		}
	}

	e.mu.RLock()
	workflow, exists := e.workflows[flowKey]
	e.mu.RUnlock()

	if !exists {
		return nil, errors.ErrWorkflowNotFound(flowKey)
	}

	// 存入缓存
	if e.cache != nil {
		e.cache.Put(flowKey, workflow)
	}

	return workflow, nil
}

// RemoveWorkflow 移除工作流
func (e *Engine) RemoveWorkflow(flowKey string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.workflows[flowKey]; !exists {
		return errors.ErrWorkflowNotFound(flowKey)
	}

	delete(e.workflows, flowKey)
	if e.cache != nil {
		e.cache.Delete(flowKey)
	}

	return nil
}

// ExecuteWorkflowWithContext 带上下文执行工作流
func (e *Engine) ExecuteWorkflowWithContext(ctx context.Context, flowKey string, workflowCtx *types.WorkflowContext) (map[string]interface{}, error) {
	// 获取工作流
	workflow, err := e.GetWorkflow(flowKey)
	if err != nil {
		return nil, err
	}

	// 检查并发限制
	select {
	case e.semaphore <- struct{}{}:
		defer func() { <-e.semaphore }()
	case <-ctx.Done():
		return nil, errors.ErrExecutionTimeout(flowKey)
	}

	// 设置超时
	timeout := time.Duration(e.config.DefaultTimeout) * time.Second
	if len(workflow.Nodes) > 0 && workflow.Nodes[0].Timeout > 0 {
		timeout = time.Duration(workflow.Nodes[0].Timeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 执行工作流
	return e.executeWorkflowInternal(ctx, workflow, workflowCtx)
}

// executeWorkflowInternal 内部执行工作流逻辑
func (e *Engine) executeWorkflowInternal(ctx context.Context, workflow *types.Workflow, workflowCtx *types.WorkflowContext) (map[string]interface{}, error) {
	// 从第一个节点开始
	if len(workflow.Nodes) == 0 {
		return nil, errors.ErrInvalidConfig("workflow has no nodes")
	}

	// 初始化 Metadata（如果为空）
	if workflowCtx.Metadata == nil {
		workflowCtx.Metadata = make(map[string]interface{})
	}

	currentNodeID := workflow.Nodes[0].ID

	for {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return nil, errors.ErrExecutionTimeout(workflow.Metadata.FlowKey)
		default:
		}

		// 查找当前节点
		var currentNode *types.Node
		for i := range workflow.Nodes {
			if workflow.Nodes[i].ID == currentNodeID {
				currentNode = &workflow.Nodes[i]
				break
			}
		}

		if currentNode == nil {
			return nil, errors.ErrNodeNotFound(currentNodeID)
		}

		// 执行节点
		executor := executors.GetNodeExecutorWithEngine(currentNode.Type, e)
		resultWorkflowCtx, err := executor.Execute(ctx, currentNode, workflowCtx, workflow)
		if err != nil {
			return nil, errors.ErrExecutionFailed(currentNodeID, err)
		}

		// 更新工作流上下文
		workflowCtx = resultWorkflowCtx

		// 检查是否结束（从工作流上下文中获取下一个节点ID）
		nextNodeID := workflowCtx.NextNodeID
		if nextNodeID == 0 {
			break
		}

		currentNodeID = nextNodeID
	}

	// 返回业务数据
	return workflowCtx.GetData(), nil
}

// ExecuteWorkflowAsync 异步执行工作流
func (e *Engine) ExecuteWorkflowAsync(ctx context.Context, flowKey string, workflowCtx *types.WorkflowContext, callback func(map[string]interface{}, error)) {
	go func() {
		result, err := e.ExecuteWorkflowWithContext(ctx, flowKey, workflowCtx)
		if callback != nil {
			callback(result, err)
		}
	}()
}

// GetWorkflowCount 获取工作流数量
func (e *Engine) GetWorkflowCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.workflows)
}

// GetConfig 获取引擎配置
func (e *Engine) GetConfig() *Config {
	return e.config
}

// GetRuleEngine 获取规则引擎
func (e *Engine) GetRuleEngine() *ruleengine.Engine {
	return e.ruleEngine
}

// ClearCache 清空缓存
func (e *Engine) ClearCache() {
	if e.cache != nil {
		e.cache.Clear()
	}
}

// Shutdown 关闭引擎
func (e *Engine) Shutdown(ctx context.Context) error {
	_ = ctx // 保留参数以符合接口定义
	var err error
	e.shutdownOnce.Do(func() {
		// 等待所有正在执行的工作流完成
		for i := 0; i < e.config.MaxConcurrent; i++ {
			e.semaphore <- struct{}{}
		}

		// 清空缓存
		if e.cache != nil {
			e.cache.Clear()
		}
	})

	return err
}

// WorkflowCache 方法

// Put 存入缓存
func (c *WorkflowCache) Put(key string, workflow *types.Workflow) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &cacheEntry{
		workflow: workflow,
		expires:  time.Now().Add(c.ttl),
	}
}

// Get 从缓存获取
func (c *WorkflowCache) Get(key string) *types.Workflow {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil
	}

	// 检查是否过期
	if time.Now().After(entry.expires) {
		delete(c.cache, key)
		return nil
	}

	return entry.workflow
}

// Delete 删除缓存
func (c *WorkflowCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

// Clear 清空缓存
func (c *WorkflowCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]*cacheEntry)
}

// Cleanup 清理过期缓存
func (c *WorkflowCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.cache {
		if now.After(entry.expires) {
			delete(c.cache, key)
		}
	}
}
