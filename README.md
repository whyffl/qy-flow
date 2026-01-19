# 工作流引擎和规则引擎

基于Golang实现的高并发、低延迟的工作流引擎和规则引擎系统，支持JSON格式DSL配置。

## 项目特性

- ✅ **高并发**: 使用goroutines和sync.WaitGroup实现并发执行
- ✅ **低延迟**: 使用sync.Pool、预编译正则表达式、sync.Map等优化手段
- ✅ **JSON格式DSL**: 简洁易用的JSON格式工作流和规则定义
- ✅ **可扩展**: 支持自定义规则类型和节点执行器
- ✅ **多种节点类型**: 支持start、switch、action、parallel、subprocess、end节点
- ✅ **多种规则类型**: 支持字符串匹配、数值比较、自定义处理器
- ✅ **工作流缓存**: 内置缓存机制提升性能
- ✅ **异步执行**: 支持同步和异步工作流执行
- ✅ **超时控制**: 支持上下文超时控制
- ✅ **优雅关闭**: 支持优雅的引擎关闭

## 项目结构

```
workflow/
├── cmd/
│   └── main.go                 # 主程序入口
├── configs/
│   ├── application.json        # 应用配置
│   ├── workflow_example.json   # 工作流示例
│   └── workflow_parallel_example.json # 并行工作流示例
├── pkg/
│   ├── constant/
│   │   └── constant.go         # 常量定义
│   ├── errors/
│   │   └── errors.go           # 错误定义
│   ├── ruleengine/
│   │   ├── types.go            # 规则类型定义
│   │   ├── string_match.go     # 字符串匹配规则
│   │   ├── numeric_compare.go  # 数值比较规则
│   │   ├── handler.go          # 处理器规则
│   │   ├── factory.go          # 规则工厂
│   │   └── engine.go           # 规则引擎
│   ├── workflowengine/
│   │   ├── types.go            # 工作流类型定义
│   │   ├── config.go           # 配置管理
│   │   ├── node_executors.go   # 节点执行器
│   │   └── engine.go           # 工作流引擎
│   └── api/
│       └── api.go              # API接口
├── docs/
│   └── 提示词.txt               # 项目需求文档
└── go.mod                      # Go模块文件
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 运行示例

```bash
go run cmd/main.go
```

## 使用说明

### 创建工作流引擎

```go
// 创建规则引擎
ruleEngine := ruleengine.NewEngine(
    ruleengine.WithTimeout(30*time.Second),
    ruleengine.WithWorkerPoolSize(50),
)

// 创建工作流引擎
wfEngine := workflowengine.NewEngine(nil, ruleEngine)

// 创建组合API
engineAPI := api.NewEngineAPI(wfEngine, ruleEngine)
```

### 从JSON加载工作流

```go
// 加载工作流定义
workflowDef, err := workflowengine.LoadWorkflowDefinition("configs/workflow_example.json")
if err != nil {
    log.Fatal(err)
}

// 创建工作流
if err := engineAPI.Workflow.CreateWorkflowFromDefinition(workflowDef); err != nil {
    log.Fatal(err)
}
```

### 执行工作流

```go
data := map[string]interface{}{
    "user_id":   "1001",
    "user_name": "Alice",
    "user_type": "vip",
    "age":       28.0,
}

result, err := engineAPI.Workflow.ExecuteWorkflow("user_approval_flow", data)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("结果: %+v\n", result)
```

### 异步执行

```go
engineAPI.Workflow.ExecuteWorkflowAsync("user_approval_flow", data, func(result map[string]interface{}, err error) {
    if err != nil {
        log.Printf("执行失败: %v", err)
    } else {
        fmt.Printf("结果: %+v\n", result)
    }
})
```

### 创建规则

```go
// 创建字符串匹配规则
err := engineAPI.Rule.CreateRule("check_vip",
    ruleengine.StringMatch,
    ruleengine.RuleConfig{
        Field:     "user_type",
        Operation: "equals",
        Value:     "vip",
    })

// 创建数值比较规则
err := engineAPI.Rule.CreateRule("check_age",
    ruleengine.NumericCompare,
    ruleengine.RuleConfig{
        Field:     "age",
        Operation: "greater_than",
        Value:     18.0,
    })
```

### 并发执行规则

```go
results, err := engineAPI.Rule.ExecuteRulesConcurrent([]string{"rule1", "rule2", "rule3"}, data)
if err != nil {
    log.Fatal(err)
}

for ruleID, result := range results {
    fmt.Printf("%s: %v\n", ruleID, result)
}
```

## DSL格式说明

### 工作流定义JSON格式

```json
{
  "metadata": {
    "id": 1,
    "flow_key": "workflow_key",
    "name": "工作流名称",
    "version": "1.0.0",
    "status": "published"
  },
  "nodes": [
    {
      "id": 1,
      "type": "start",
      "name": "开始",
      "next": 2
    },
    {
      "id": 2,
      "type": "switch",
      "name": "分支节点",
      "cases": [
        {
          "condition_rule_id": "rule_id",
          "next": 3
        }
      ],
      "default": {
        "next": 4
      }
    },
    {
      "id": 3,
      "type": "action",
      "name": "动作节点",
      "handler": "handler_name",
      "next": 5
    },
    {
      "id": 4,
      "type": "parallel",
      "name": "并行节点",
      "parallel_nodes": [6, 7, 8],
      "next": 9
    },
    {
      "id": 5,
      "type": "end",
      "name": "结束"
    }
  ],
  "rules": [
    {
      "rule_id": "rule_id",
      "type": "string_match",
      "config": {
        "field": "field_name",
        "operation": "equals",
        "value": "expected_value"
      }
    }
  ]
}
```

### 节点类型

| 类型 | 说明 | 字段 |
|------|------|------|
| start | 开始节点 | next |
| switch | 分支节点 | cases, default |
| action | 动作节点 | handler, next |
| parallel | 并行节点 | parallel_nodes, next |
| subprocess | 子流程节点 | subprocess_key, next |
| end | 结束节点 | - |

### 规则类型

| 类型 | 说明 |
|------|------|
| string_match | 字符串匹配规则 |
| numeric_compare | 数值比较规则 |
| handler | 自定义处理器规则 |

### 字符串操作符

- equals: 等于
- not_equals: 不等于
- contains: 包含
- not_contains: 不包含
- starts_with: 以...开头
- ends_with: 以...结尾
- regex: 正则匹配

### 数值操作符

- equal: 等于
- not_equal: 不等于
- greater_than: 大于
- greater_or_equal: 大于等于
- less_than: 小于
- less_or_equal: 小于等于

## 配置说明

### 应用配置 (application.json)

```json
{
  "max_concurrent": 1000,           // 最大并发工作流数量
  "default_timeout": 30,             // 默认超时时间（秒）
  "worker_pool_size": 100,           // 工作协程池大小
  "enable_cache": true,             // 是否启用缓存
  "cache_ttl": 3600,                // 缓存过期时间（秒）
  "shutdown_timeout": "30s"          // 关闭超时时间
}
```

## 性能优化

1. **正则表达式预编译**: 字符串匹配规则中的正则表达式会被预编译并缓存
2. **工作流缓存**: 工作流定义会被缓存，减少重复加载开销
3. **协程池**: 使用固定大小的协程池限制并发数量
4. **sync.Map**: 使用并发安全的map存储规则和工作流
5. **sync.Pool**: 使用对象池复用对象，减少GC压力

## 扩展开发

### 自定义规则类型

```go
// 实现Rule接口
type CustomRule struct {
    Config ruleengine.RuleConfig
}

func (r *CustomRule) Execute(data map[string]interface{}) bool {
    // 自定义执行逻辑
    return true
}

func (r *CustomRule) GetType() string {
    return "custom_type"
}

// 注册到工厂
ruleengine.RegisterCustomRule("custom_type", func(config ruleengine.RuleConfig) ruleengine.Rule {
    return &CustomRule{Config: config}
})
```

### 自定义节点执行器

```go
// 实现NodeExecutor接口
type CustomNodeExecutor struct{}

func (e *CustomNodeExecutor) Execute(ctx context.Context, node *workflowengine.Node, data map[string]interface{}, workflow *workflowengine.Workflow) (map[string]interface{}, error) {
    // 自定义执行逻辑
    data["result"] = "custom_result"
    data["_next_node_id"] = node.Next
    return data, nil
}

// 注册到工作流引擎
workflowengine.RegisterNodeExecutor("custom_node", &CustomNodeExecutor{})
```

## 注意事项

1. 节点ID必须唯一且连续
2. 节点的Next字段指向的节点ID必须存在
3. 规则ID在工作流中必须唯一
4. Switch节点必须有cases或default之一
5. Parallel节点必须有parallel_nodes列表
6. 建议在生产环境中配置适当的超时时间

## 许可证

MIT License
