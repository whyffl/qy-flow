package executors

import (
	"context"
	"sync"
	"workflow/pkg/errors"
	"workflow/pkg/workflowengine/types"
	"workflow/pkg/workflowengine/utils"
)

// ParallelNodeExecutor 并行节点执行器
type ParallelNodeExecutor struct{}

// Execute 执行并行节点
func (e *ParallelNodeExecutor) Execute(ctx context.Context, node *types.Node, workflowCtx *types.WorkflowContext, workflow *types.Workflow) (*types.WorkflowContext, error) {
	if len(node.ParallelNodes) == 0 {
		workflowCtx.SetNextNode(node.Next)
		return workflowCtx, nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(node.ParallelNodes))
	results := make(map[int]*types.WorkflowContext)

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

			// 执行子节点，创建新的工作流上下文
			childWorkflowCtx := &types.WorkflowContext{
				Data: utils.DeepCopyMap(workflowCtx.GetData()),
			}
			executor := GetNodeExecutor(childNode.Type)
			result, err := executor.Execute(ctx, childNode, childWorkflowCtx, workflow)

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

	// 合并所有结果到主上下文的业务数据中
	for _, result := range results {
		for key, value := range result.GetData() {
			workflowCtx.GetData()[key] = value
		}
	}

	// 设置下一个节点（在上下文中，而不是修改业务数据）
	workflowCtx.SetNextNode(node.Next)
	return workflowCtx, nil
}
