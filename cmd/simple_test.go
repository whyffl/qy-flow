package main

import (
	"fmt"
	"log"
	"workflow/pkg/constant"

	"workflow/pkg/ruleengine"
)

func main() {
	fmt.Println("=== 规则引擎简单测试 ===")
	fmt.Println()

	// 创建规则引擎（使用默认配置）
	fmt.Println("1. 创建规则引擎...")
	ruleEngine := ruleengine.NewEngine(nil)
	fmt.Println("   ✓ 规则引擎创建成功")
	fmt.Println()

	// 创建字符串匹配规则
	fmt.Println("2. 创建字符串匹配规则...")
	err = ruleEngine.AddRuleFromDefinition("check_vip", ruleengine.RuleDefinition{
		RuleID: "check_vip",
		Type:   constant.StringMatch,
		Config: ruleengine.RuleConfig{
			Field:     "user_type",
			Operation: "equals",
			Value:     "vip",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("   ✓ 规则创建成功")
	fmt.Println()

	// 创建数值比较规则
	fmt.Println("3. 创建数值比较规则...")
	err = ruleEngine.AddRuleFromDefinition("check_age", ruleengine.RuleDefinition{
		RuleID: "check_age",
		Type:   constant.NumericCompare,
		Config: ruleengine.RuleConfig{
			Field:     "age",
			Operation: "greater_than",
			Value:     18.0,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("   ✓ 规则创建成功")
	fmt.Println()

	// 测试规则执行
	fmt.Println("4. 测试规则执行...")
	data := map[string]interface{}{
		"user_type": "vip",
		"age":       25.0,
	}

	result1, err := ruleEngine.ExecuteRule("check_vip", data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   - check_vip 结果: %v\n", result1)

	result2, err := ruleEngine.ExecuteRule("check_age", data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   - check_age 结果: %v\n\n", result2)

	// 并发执行规则
	fmt.Println("5. 并发执行规则...")
	results, err := ruleEngine.ExecuteRulesConcurrent([]string{"check_vip", "check_age"}, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("   ✓ 并发执行成功")
	for ruleID, result := range results {
		fmt.Printf("   - %s: %v\n", ruleID, result)
	}

	fmt.Println("\n=== 测试完成 ===")
}
