package types

// WorkflowContext 工作流执行上下文
type WorkflowContext struct {
	NextNodeID int                    // 下一个节点ID（0表示结束）
	Data       map[string]interface{} // 业务数据
	Metadata   map[string]interface{} // 执行元数据
}

// SetNextNode 设置下一个节点
func (c *WorkflowContext) SetNextNode(nodeID int) {
	c.NextNodeID = nodeID
}

// SetData 设置业务数据
func (c *WorkflowContext) SetData(data map[string]interface{}) {
	c.Data = data
}

// GetData 获取业务数据
func (c *WorkflowContext) GetData() map[string]interface{} {
	return c.Data
}

// SetMetadata 设置元数据
func (c *WorkflowContext) SetMetadata(key string, value interface{}) {
	if c.Metadata == nil {
		c.Metadata = make(map[string]interface{})
	}
	c.Metadata[key] = value
}

// GetMetadata 获取元数据
func (c *WorkflowContext) GetMetadata(key string) (interface{}, bool) {
	if c.Metadata == nil {
		return nil, false
	}
	val, ok := c.Metadata[key]
	return val, ok
}
