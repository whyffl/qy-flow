package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"workflow/internal/config"
	"workflow/internal/workflow"
	"workflow/pkg/constant"

	"workflow/pkg/api"
	"workflow/pkg/ruleengine"
	"workflow/pkg/workflowengine"
	"workflow/pkg/workflowengine/types"
)

func main() {
	fmt.Println("=== 工作流引擎和规则引擎演示 ===")
	fmt.Println()

	// 加载配置文件
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
	fmt.Printf("配置加载成功: %s v%s\n\n", cfg.App.Name, cfg.App.Version)

	// 1. 创建规则注册表（可选）
	fmt.Println("1. 创建规则注册表...")
	registry := ruleengine.NewRuleRegistry()

	// 注册自定义 handler 示例
	fmt.Println("2. 注册自定义 handler...")
	if err := registry.RegisterHandler("custom_approve", func(data map[string]interface{}) bool {
		// 自定义审批逻辑：金额小于1000自动通过
		if amount, ok := data["amount"].(float64); ok {
			return amount < 1000
		}
		return false
	}); err != nil {
		log.Printf("注册 handler 失败: %v", err)
	} else {
		fmt.Println("   ✓ custom_approve handler 注册成功")
	}

	if err := registry.RegisterHandler("risk_check", func(data map[string]interface{}) bool {
		// 风险检查：用户类型为 low_risk 时通过
		if riskLevel, ok := data["risk_level"].(string); ok {
			return riskLevel == "low_risk"
		}
		return false
	}); err != nil {
		log.Printf("注册 handler 失败: %v", err)
	} else {
		fmt.Println("   ✓ risk_check handler 注册成功")
	}
	fmt.Println()

	// 3. 创建规则引擎（传入配置和注册表）
	fmt.Println("3. 创建规则引擎...")
	ruleEngineConfig := cfg.RuleEngine.ToEngineConfig()
	ruleEngine := ruleengine.NewEngine(ruleEngineConfig, registry)
	fmt.Println("   ✓ 规则引擎创建成功")
	fmt.Printf("   - 工作池大小: %d\n", ruleEngineConfig.WorkerPoolSize)
	fmt.Printf("   - 正则池大小: %d\n", ruleEngineConfig.RegexPoolSize)
	fmt.Println()

	// 4. 创建工作流引擎（使用配置）
	fmt.Println("4. 创建工作流引擎...")
	wfEngineConfig := cfg.WorkflowEngine.ToEngineConfig()
	wfEngine := workflowengine.NewEngine(wfEngineConfig, ruleEngine)
	fmt.Println("   ✓ 工作流引擎创建成功")
	fmt.Printf("   - 最大并发数: %d\n", wfEngineConfig.MaxConcurrent)
	fmt.Printf("   - 默认超时: %d秒\n", wfEngineConfig.DefaultTimeout)
	fmt.Printf("   - 工作池大小: %d\n", wfEngineConfig.WorkerPoolSize)
	fmt.Printf("   - 缓存启用: %v\n", wfEngineConfig.EnableCache)
	fmt.Println()

	// 5. 创建组合API
	engineAPI := api.NewEngineAPI(wfEngine, ruleEngine)

	// 6. 从JSON文件加载工作流定义
	fmt.Println("6. 从JSON文件加载工作流定义...")
	workflowDef, err := workflow.LoadWorkflowDefinition("configs/workflow_example.json")
	if err != nil {
		log.Fatalf("加载工作流定义失败: %v", err)
	}
	fmt.Println("   ✓ 工作流定义加载成功")
	fmt.Printf("   - 流程名称: %s\n", workflowDef.Metadata.Name)
	fmt.Printf("   - 版本: %s\n", workflowDef.Metadata.Version)
	fmt.Printf("   - 节点数: %d\n", len(workflowDef.Nodes))
	fmt.Printf("   - 规则数: %d\n\n", len(workflowDef.Rules))

	// 7. 创建工作流
	fmt.Println("7. 创建工作流...")
	if err := engineAPI.Workflow.CreateWorkflowFromDefinition(workflowDef); err != nil {
		log.Fatalf("创建工作流失败: %v", err)
	}
	fmt.Println("   ✓ 工作流创建成功")
	fmt.Println()

	// 8. 演示VIP用户流程
	fmt.Println("8. 执行VIP用户审批流程...")
	vipData := map[string]interface{}{
		"user_id":   "1001",
		"user_name": "Alice",
		"user_type": "vip",
		"age":       28.0,
	}

	ctx := context.Background()
	vipWorkflowCtx := &types.WorkflowContext{
		Data:     vipData,
		Metadata: make(map[string]interface{}),
	}

	result, err := engineAPI.Workflow.ExecuteWorkflowWithContext(ctx, "user_approval_flow", vipWorkflowCtx)
	if err != nil {
		log.Printf("执行失败: %v", err)
	} else {
		fmt.Println("   ✓ 执行成功")
		fmt.Printf("   - 结果: %+v\n\n", result)
	}

	// 9. 演示新用户流程
	fmt.Println("9. 执行新用户审批流程...")
	newData := map[string]interface{}{
		"user_id":   "1002",
		"user_name": "Bob",
		"user_type": "new",
		"age":       22.0,
	}

	newWorkflowCtx := &types.WorkflowContext{
		Data:     newData,
		Metadata: make(map[string]interface{}),
	}

	result, err = engineAPI.Workflow.ExecuteWorkflowWithContext(ctx, "user_approval_flow", newWorkflowCtx)
	if err != nil {
		log.Printf("执行失败: %v", err)
	} else {
		fmt.Println("   ✓ 执行成功")
		fmt.Printf("   - 结果: %+v\n\n", result)
	}

	// 10. 演示普通用户流程
	fmt.Println("10. 执行普通用户审批流程...")
	normalData := map[string]interface{}{
		"user_id":   "1003",
		"user_name": "Charlie",
		"user_type": "normal",
		"age":       35.0,
	}

	normalWorkflowCtx := &types.WorkflowContext{
		Data:     normalData,
		Metadata: make(map[string]interface{}),
	}

	result, err = engineAPI.Workflow.ExecuteWorkflowWithContext(ctx, "user_approval_flow", normalWorkflowCtx)
	if err != nil {
		log.Printf("执行失败: %v", err)
	} else {
		fmt.Println("   ✓ 执行成功")
		fmt.Printf("   - 结果: %+v\n\n", result)
	}

	// 11. 演示异步执行
	fmt.Println("11. 异步执行工作流...")
	asyncWorkflowCtx := &types.WorkflowContext{
		Data:     vipData,
		Metadata: make(map[string]interface{}),
	}
	engineAPI.Workflow.ExecuteWorkflowAsync(ctx, "user_approval_flow", asyncWorkflowCtx, func(result map[string]interface{}, err error) {
		if err != nil {
			log.Printf("异步执行失败: %v", err)
		} else {
			fmt.Println("   ✓ 异步执行成功")
			fmt.Printf("   - 结果: %+v\n", result)
		}
	})
	time.Sleep(100 * time.Millisecond) // 等待异步执行完成
	fmt.Println()

	// 12. 演示规则引擎
	fmt.Println("12. 演示规则引擎...")
	fmt.Println("   a. 创建字符串匹配规则...")
	if err := engineAPI.Rule.CreateRule("check_age",
		constant.StringMatch,
		ruleengine.RuleConfig{
			Field:     "status",
			Operation: "equals",
			Value:     "active",
		}); err != nil {
		log.Printf("创建规则失败: %v", err)
	} else {
		fmt.Println("      ✓ 规则创建成功")
	}

	fmt.Println("   b. 执行规则...")
	ruleData := map[string]interface{}{
		"status": "active",
		"name":   "test",
	}
	ruleResult, err := engineAPI.Rule.ExecuteRule("check_age", ruleData)
	if err != nil {
		log.Printf("执行规则失败: %v", err)
	} else {
		fmt.Printf("      ✓ 规则执行结果: %v\n\n", ruleResult)
	}

	// 13. 并行执行规则
	fmt.Println("13. 并发执行多个规则...")
	if err := engineAPI.Rule.CreateRule("rule1", constant.StringMatch, ruleengine.RuleConfig{
		Field: "field1", Operation: "equals", Value: "value1",
	}); err != nil {
		log.Printf("创建规则rule1失败: %v", err)
	}
	if err := engineAPI.Rule.CreateRule("rule2", constant.NumericCompare, ruleengine.RuleConfig{
		Field: "field2", Operation: "greater_than", Value: 10.0,
	}); err != nil {
		log.Printf("创建规则rule2失败: %v", err)
	}
	if err := engineAPI.Rule.CreateRule("rule3", constant.StringMatch, ruleengine.RuleConfig{
		Field: "field3", Operation: "contains", Value: "substring",
	}); err != nil {
		log.Printf("创建规则rule3失败: %v", err)
	}

	multiRuleData := map[string]interface{}{
		"field1": "value1",
		"field2": 15.0,
		"field3": "this is a substring example",
	}

	results, err := engineAPI.Rule.ExecuteRulesConcurrent([]string{"rule1", "rule2", "rule3"}, multiRuleData)
	if err != nil {
		log.Printf("并发执行规则失败: %v", err)
	} else {
		fmt.Println("   ✓ 并发执行成功")
		for ruleID, result := range results {
			fmt.Printf("   - %s: %v\n", ruleID, result)
		}
	}
	fmt.Println()

	// 14. 演示并行节点工作流
	fmt.Println("14. 加载并执行并行节点工作流...")
	parallelWfDef, err := workflow.LoadWorkflowDefinition("configs/workflow_parallel_example.json")
	if err != nil {
		log.Printf("加载并行工作流失败: %v", err)
	} else {
		if err := engineAPI.Workflow.CreateWorkflowFromDefinition(parallelWfDef); err != nil {
			log.Printf("创建并行工作流失败: %v", err)
		} else {
			fmt.Println("   ✓ 并行工作流创建成功")

			parallelData := map[string]interface{}{
				"task_id": "TASK-001",
			}
			parallelWorkflowCtx := &types.WorkflowContext{
				Data:     parallelData,
				Metadata: make(map[string]interface{}),
			}
			parallelResult, err := engineAPI.Workflow.ExecuteWorkflowWithContext(ctx, "parallel_task_flow", parallelWorkflowCtx)
			if err != nil {
				log.Printf("执行并行工作流失败: %v", err)
			} else {
				fmt.Println("   ✓ 并行工作流执行成功")
				fmt.Printf("   - 结果: %+v\n\n", parallelResult)
			}
		}
	}

	// 15. 获取统计信息
	fmt.Println("15. 获取引擎统计信息...")
	wfStats := engineAPI.Workflow.GetStats()
	ruleStats := engineAPI.Rule.GetStats()
	fmt.Printf("   - 工作流数量: %d\n", wfStats.WorkflowCount)
	fmt.Printf("   - 规则数量: %d\n", ruleStats.RuleCount)
	fmt.Println()

	// 16. 带超时的执行
	fmt.Println("16. 带超时的执行...")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用 API 执行工作流
	timeoutWorkflowCtx := &types.WorkflowContext{
		Data:     vipData,
		Metadata: make(map[string]interface{}),
	}
	_, err = engineAPI.Workflow.ExecuteWorkflowWithContext(timeoutCtx, "user_approval_flow", timeoutWorkflowCtx)
	if err != nil {
		fmt.Printf("   - 执行超时或失败: %v\n", err)
	} else {
		fmt.Println("   ✓ 执行成功（在超时时间内）")
	}
	fmt.Println()

	// 17. 优雅关闭
	fmt.Println("17. 优雅关闭引擎...")
	if err := wfEngine.Shutdown(context.Background()); err != nil {
		log.Printf("关闭引擎失败: %v", err)
	} else {
		fmt.Println("   ✓ 引擎关闭成功")
	}

	fmt.Println("\n=== 演示完成 ===")
}
