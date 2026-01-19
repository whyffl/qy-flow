package workflow

import (
	"encoding/json"
	"os"
	"workflow/pkg/errors"
	"workflow/pkg/workflowengine/types"
)

// LoadWorkflowDefinition 从JSON文件加载工作流定义
func LoadWorkflowDefinition(filepath string) (*types.WorkflowDefinition, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, errors.ErrInvalidConfig(err.Error())
	}

	var def types.WorkflowDefinition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, errors.ErrInvalidConfig(err.Error())
	}

	return &def, nil
}
