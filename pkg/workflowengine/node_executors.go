package workflowengine

import (
	"context"
	"fmt"
	"sync"
	"workflow/pkg/constant"
	"workflow/pkg/errors"
)

// NodeExecutor 节点执行器接口
type NodeExecutor interface {
	Execute(ctx context.Context, node *Node, data map[string]interface{}, workflow *Workflow) (map[string]interface{}, error)
}

// StartNodeExecutor 开始节点执行器
type StartNodeExecutor struct{}

// Execute 执行开始节点
func (e *StartNodeExecutor) Execute(ctx context.Context, node *Node, data map[string]interface{}, workflow *Workflow) (map[string]interface{}, error) {
	// 参数声明但未使用，符合接口定义要求
	_ = ctx
	_ = node
	_ = workflow

	// 开始节点通常不做任何处理，只是传递数据
	return data, nil
}

// SwitchNodeExecutor 分支节点执行器
type SwitchNodeExecutor struct{}

// Execute 执行分支节点
func (e *SwitchNodeExecutor) Execute(ctx context.Context, node *Node, data map[string]interface{}, workflow *Workflow) (map[string]interface{}, error) {
	_ = ctx // 参数声明但未使用，符合接口定义要求
	// 遍历所有条件分支，找到第一个匹配的规则
	for _, caseItem := range node.Cases {
		rule, exists := workflow.Rules[caseItem.ConditionRuleID]
		if !exists {
			return nil, errors.ErrFieldNotFound(caseItem.ConditionRuleID)
		}

		if rule.Execute(data) {
			// 设置下一个节点
			data["_next_node_id"] = caseItem.Next
			return data, nil
		}
	}

	// 如果没有匹配的条件，使用默认分支
	if node.Default != nil {
		data["_next_node_id"] = node.Default.Next
		return data, nil
	}

	// 没有匹配的条件也没有默认分支，工作流结束
	data["_next_node_id"] = 0
	return data, nil
}

// ActionNodeExecutor 动作节点执行器
type ActionNodeExecutor struct{}

// Execute 执行动作节点
func (e *ActionNodeExecutor) Execute(ctx context.Context, node *Node, data map[string]interface{}, workflow *Workflow) (map[string]interface{}, error) {
	_ = ctx
	_ = workflow

	if node.Handler == "" {
		return data, nil
	}

	// TODO: 实现handler注册和执行机制
	// 这里可以扩展为调用外部handler
	fmt.Printf("Executing action handler: %s\n", node.Handler)

	// 设置下一个节点
	data["_next_node_id"] = node.Next
	return data, nil
}

// ParallelNodeExecutor 并行节点执行器
type ParallelNodeExecutor struct{}

// Execute 执行并行节点
func (e *ParallelNodeExecutor) Execute(ctx context.Context, node *Node, data map[string]interface{}, workflow *Workflow) (map[string]interface{}, error) {
	if len(node.ParallelNodes) == 0 {
		data["_next_node_id"] = node.Next
		return data, nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(node.ParallelNodes))
	results := make(map[int]map[string]interface{})

	// 并行执行所有子节点
	for _, childNodeID := range node.ParallelNodes {
		wg.Add(1)

		go func(nodeID int) {
			defer wg.Done()

			// 查找子节点
			var childNode *Node
			for i := range workflow.Nodes {
				if workflow.Nodes[i].ID == nodeID {
					childNode = &workflow.Nodes[i]
					break
				}
			}

			if childNode == nil {
				errChan <- errors.ErrNodeNotFound(nodeID)
				return
			}

			// 执行子节点
			resultData := deepCopyMap(data)
			executor := GetNodeExecutor(childNode.Type)
			result, err := executor.Execute(ctx, childNode, resultData, workflow)

			mu.Lock()
			if err != nil {
				errChan <- err
			} else {
				results[nodeID] = result
			}
			mu.Unlock()
		}(childNodeID)
	}

	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		if err != nil {
			return nil, errors.ErrParallelFailed(err)
		}
	}

	// 合并所有结果（合并非保留字段）
	for _, result := range results {
		for key, value := range result {
			if key != "_next_node_id" {
				data[key] = value
			}
		}
	}

	data["_next_node_id"] = node.Next
	return data, nil
}

// SubprocessNodeExecutor 子流程节点执行器
type SubprocessNodeExecutor struct {
	engine *Engine
}

// Execute 执行子流程节点
func (e *SubprocessNodeExecutor) Execute(ctx context.Context, node *Node, data map[string]interface{}, workflow *Workflow) (map[string]interface{}, error) {
	_ = workflow

	if node.SubprocessKey == "" {
		return nil, errors.ErrInvalidConfig("subprocess key is required")
	}

	// 检查子流程是否存在
	if _, err := e.engine.GetWorkflow(node.SubprocessKey); err != nil {
		return nil, err
	}

	// 执行子流程
	subResult, err := e.engine.ExecuteWorkflowWithContext(ctx, node.SubprocessKey, data)
	if err != nil {
		return nil, err
	}

	// 合并子流程结果
	for key, value := range subResult {
		data[key] = value
	}

	data["_next_node_id"] = node.Next
	return data, nil
}

// EndNodeExecutor 结束节点执行器
type EndNodeExecutor struct{}

// Execute 执行结束节点
func (e *EndNodeExecutor) Execute(ctx context.Context, node *Node, data map[string]interface{}, workflow *Workflow) (map[string]interface{}, error) {
	_ = ctx
	_ = node
	_ = workflow

	data["_next_node_id"] = 0
	return data, nil
}

// GetNodeExecutor 根据节点类型获取对应的执行器
func GetNodeExecutor(nodeType string) NodeExecutor {
	switch nodeType {
	case constant.Start:
		return &StartNodeExecutor{}
	case constant.Switch:
		return &SwitchNodeExecutor{}
	case constant.Action:
		return &ActionNodeExecutor{}
	case constant.Parallel:
		return &ParallelNodeExecutor{}
	case constant.Subprocess:
		return &SubprocessNodeExecutor{}
	case constant.End:
		return &EndNodeExecutor{}
	default:
		return &ActionNodeExecutor{} // 默认使用action执行器
	}
}

// deepCopyMap 深拷贝map
func deepCopyMap(original map[string]interface{}) map[string]interface{} {
	copiedMap := make(map[string]interface{})
	for key, value := range original {
		copiedMap[key] = value
	}
	return copiedMap
}
