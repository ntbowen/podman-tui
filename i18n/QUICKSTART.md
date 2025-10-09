# i18n 快速入门

## 5 分钟上手指南

### 1️⃣ 导入包

```go
import "github.com/containers/podman-tui/i18n"
```

### 2️⃣ 翻译文本

```go
// 之前
button := "Cancel"

// 之后
button := i18n.T("Cancel")  // 自动翻译为 "取消"（中文环境）
```

### 3️⃣ 处理宽度

```go
// ❌ 错误 - 中文会出问题
width := len("删除")  // 返回 6（字节数）

// ✅ 正确 - 使用 i18n 函数
width := i18n.GetDisplayWidth("删除")  // 返回 4（显示宽度）
```

## 常用函数速查

### 翻译

```go
i18n.T("Cancel")                    // 翻译字符串（默认英文）
i18n.SetLanguage("zh_cn")           // 设置为中文
i18n.GetCurrentLanguage()           // 获取当前语言 "zh_cn"
i18n.GetSupportedLanguages()        // ["en", "zh_cn", "zh_tw", "ja", "ko", "ru"]
i18n.GetLanguageName("zh_cn")       // "简体中文"
```

### 宽度处理

```go
i18n.GetDisplayWidth("删除")                    // 获取显示宽度: 4
i18n.PadToWidth("删除", 20)                     // 填充到20: "删除                "
i18n.CalcLabelWidth("容器名称", 2)              // 标签宽度: 10
i18n.FormatLabelWithWidth("删除", 10, ":")     // 格式化: "      删除:"
i18n.AdjustWidthForLocale(20)                  // 自动调整: 30 (中文)
i18n.TruncateToWidth("很长的文本...", 10)      // 截断: "很长的..."
i18n.AlignLabels([]string{"名称", "状态"}, ":") // 对齐标签
```

## 实际例子

### 按钮

```go
// 之前
cancelBtn := "Cancel"
okBtn := "OK"

// 之后
cancelBtn := i18n.T("Cancel")  // "取消"
okBtn := i18n.T("OK")          // "确定"
```

### 表单

```go
// 之前
form.AddInputField("Name:", "", 20, nil, nil)

// 之后
form.AddInputField(i18n.T("Name")+":", "", 20, nil, nil)
```

### 表格对齐

```go
// 之前（错误）
headers := []string{"名称", "状态", "镜像"}
for _, h := range headers {
    width := len(h) + 2  // ❌ 中文会错位
}

// 之后（正确）
headers := []string{i18n.T("Name"), i18n.T("Status"), i18n.T("Image")}
for _, h := range headers {
    width := i18n.GetDisplayWidth(h) + 2  // ✅ 正确对齐
}
```

### 错误消息

```go
// 之前
return fmt.Errorf("connection failed: %w", err)

// 之后
return fmt.Errorf("%s: %w", i18n.T("connection failed"), err)
```

## 支持的语言

| 代码 | 语言 | 翻译条目 |
|------|------|---------|
| en | English | 默认 |
| zh_cn | 简体中文 | 761 |
| zh_tw | 繁體中文 | 761 |
| ja | 日本語 | 761 |
| ko | 한국어 | 761 |
| ru | Русский | 761 |

## 设置语言

在应用启动时设置语言：

```go
// 在 main() 函数中
import "github.com/containers/podman-tui/i18n"

func main() {
    // 根据用户配置或系统环境设置语言
    i18n.SetLanguage("zh_cn")  // 设置为中文
    
    // 或者从配置文件读取
    lang := config.GetLanguage()  // 例如: "ja"
    i18n.SetLanguage(lang)
    
    // 应用的其他初始化代码...
}
```

## 运行示例

```bash
# 运行完整示例程序
go run ./i18n/example/main.go

# 运行测试
go test ./i18n/

# 检查翻译完整性
go run ./i18n/tools/check_translations.go
```

## 常见错误

### ❌ 使用 len() 计算宽度

```go
width := len("删除") + 1  // 返回 7，但显示宽度是 5
```

### ✅ 使用 i18n 函数

```go
width := i18n.CalcLabelWidth("删除", 1)  // 返回 5
```

### ❌ 硬编码字符串

```go
button := "Delete"  // 不会翻译
```

### ✅ 使用翻译函数

```go
button := i18n.T("Delete")  // 自动翻译
```

## 需要帮助？

- 📖 详细文档: [README.md](README.md)
- 🔧 集成指南: [INTEGRATION.md](INTEGRATION.md)
- 📊 项目总结: [SUMMARY.md](SUMMARY.md)
- 💻 示例代码: [example/main.go](example/main.go)
- 🧪 单元测试: [i18n_test.go](i18n_test.go)

## 下一步

1. 在你的代码中导入 `i18n` 包
2. 用 `i18n.T()` 包装所有用户可见的字符串
3. 用 `i18n.GetDisplayWidth()` 替换 `len()` 计算宽度
4. 在不同语言环境下测试

---

**提示**: 翻译缺失时会自动回退到英文原文，所以可以放心使用！
