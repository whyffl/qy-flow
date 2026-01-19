package executors

import (
	"context"
	"sync"
	"workflow/pkg/errors"
	"workflow/pkg/workflowengine/types"
)

// ParallelNodeExecutor 并行节点执行器
type ParallelNodeExecutor struct{}

// Execute 执行并行节点
func (e *ParallelNodeExecutor) Execute(ctx context.Context, node *types.Node, data map[string]interface{}, workflow *types.Workflow) (map[string]interface{}, error) {
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
			var childNode *types.Node
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

// deepCopyMap 深拷贝map
func deepCopyMap(original map[string]interface{}) map[string]interface{} {
	copiedMap := make(map[string]interface{})
	for key, value := range original {
		copiedMap[key] = value
	}
	return copiedMap
}
