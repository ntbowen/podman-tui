# i18n 框架更新日志

## 2025-10-06 - 简化设计

### 移除的功能
- ❌ 移除了自动语言检测功能 (`DetectSystemLanguage()`)
- ❌ 移除了从环境变量自动读取语言的逻辑
- ❌ 移除了 `init()` 函数中的自动初始化

### 原因
根据项目实际需求，自动语言检测功能实现复杂度较高，且不符合实际使用场景。简化为手动设置语言的方式更加清晰明了。

### 当前设计
- ✅ 默认语言为英文 (`en`)
- ✅ 通过 `SetLanguage()` 手动设置语言
- ✅ 应用启动时根据用户配置或需求设置语言

### 代码变更

#### i18n.go
```go
// 之前
var currentLanguage string

func init() {
    currentLanguage = DetectSystemLanguage()
}

func DetectSystemLanguage() string {
    // 复杂的环境变量检测逻辑...
}

// 之后
var currentLanguage = "en"  // 默认英文

// 移除了 init() 和 DetectSystemLanguage()
```

#### i18n_test.go
```go
// 移除了 TestDetectSystemLanguage 测试
// 保留了其他所有测试
```

### 使用方式

#### 之前（自动检测）
```go
// 自动从环境变量检测语言
// 无需手动设置
```

#### 之后（手动设置）
```go
// 在应用启动时设置语言
func main() {
    // 方式1: 从配置文件读取
    lang := config.GetLanguage()
    i18n.SetLanguage(lang)
    
    // 方式2: 根据用户选择
    i18n.SetLanguage("zh_cn")
    
    // 方式3: 从环境变量读取（如果需要）
    if lang := os.Getenv("PODMAN_TUI_LANG"); lang != "" {
        i18n.SetLanguage(lang)
    }
    
    // 应用的其他代码...
}
```

### 文档更新
- ✅ 更新了 README.md
- ✅ 更新了 QUICKSTART.md
- ✅ 更新了 SUMMARY.md
- ✅ 更新了 CHECKLIST.md
- ✅ 更新了 INDEX.md
- ✅ 更新了 example/main.go

### 测试结果
```
=== RUN   TestT
--- PASS: TestT (0.00s)
=== RUN   TestGetSupportedLanguages
--- PASS: TestGetSupportedLanguages (0.00s)
=== RUN   TestGetLanguageName
--- PASS: TestGetLanguageName (0.00s)
=== RUN   TestSetLanguage
--- PASS: TestSetLanguage (0.00s)
PASS
ok      github.com/containers/podman-tui/i18n   0.529s
```

### 优势
1. **更简单** - 代码更少，逻辑更清晰
2. **更可控** - 应用完全控制语言设置
3. **更灵活** - 可以从任何来源设置语言（配置文件、命令行参数、用户选择等）
4. **更可靠** - 减少了环境变量相关的潜在问题

### 迁移指南

如果你之前依赖自动语言检测，需要做以下修改：

```go
// 在应用的 main() 函数或初始化代码中添加：
import (
    "os"
    "strings"
    "github.com/containers/podman-tui/i18n"
)

func initLanguage() {
    // 如果需要从环境变量检测，可以自己实现
    lang := os.Getenv("LANG")
    if lang != "" {
        lang = strings.ToLower(strings.Split(lang, ".")[0])
        
        // 映射到支持的语言
        switch {
        case strings.HasPrefix(lang, "zh_cn"):
            i18n.SetLanguage("zh_cn")
        case strings.HasPrefix(lang, "zh_tw"):
            i18n.SetLanguage("zh_tw")
        case strings.HasPrefix(lang, "ja"):
            i18n.SetLanguage("ja")
        case strings.HasPrefix(lang, "ko"):
            i18n.SetLanguage("ko")
        case strings.HasPrefix(lang, "ru"):
            i18n.SetLanguage("ru")
        default:
            i18n.SetLanguage("en")
        }
    }
}

func main() {
    initLanguage()
    // 应用的其他代码...
}
```

### 影响范围
- ✅ 核心功能不受影响
- ✅ 所有翻译功能正常工作
- ✅ 字符宽度处理功能正常
- ✅ 所有测试通过
- ✅ 文档已同步更新

---

**总结**: 这次更新简化了 i18n 框架的设计，移除了自动语言检测功能，改为手动设置语言。这使得框架更加简单、可控和灵活。
