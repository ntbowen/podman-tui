# i18n 框架集成指南

本文档说明如何在 podman-tui 的现有代码中集成 i18n 框架。

## 快速开始

### 1. 导入 i18n 包

在需要使用翻译的文件中导入：

```go
import "github.com/containers/podman-tui/i18n"
```

### 2. 替换硬编码字符串

**之前：**
```go
button := "Cancel"
label := "Container Name:"
```

**之后：**
```go
button := i18n.T("Cancel")
label := i18n.T("Container Name") + ":"
```

### 3. 处理字符宽度

**之前（错误）：**
```go
labelWidth := len("容器名称") + 2  // 返回 14，但实际显示宽度是 10
```

**之后（正确）：**
```go
labelWidth := i18n.CalcLabelWidth("容器名称", 2)  // 返回 10
```

## 常见使用场景

### 场景 1: 按钮和菜单项

```go
// UI 按钮
cancelButton := i18n.T("Cancel")
okButton := i18n.T("OK")
helpButton := i18n.T("HELP")

// 菜单项
menuItems := []string{
    i18n.T("add connection"),
    i18n.T("remove connection"),
    i18n.T("disk usage"),
    i18n.T("system info"),
}
```

### 场景 2: 表单标签

```go
// 创建表单标签（自动处理宽度）
labels := []string{
    i18n.T("Name"),
    i18n.T("Status"),
    i18n.T("Image"),
    i18n.T("Created"),
}

// 对齐标签到相同宽度
alignedLabels, maxWidth := i18n.AlignLabels(labels, ":")

// 使用对齐后的标签
for i, label := range alignedLabels {
    fmt.Printf("%s %s\n", label, values[i])
}
```

### 场景 3: 表格列标题

```go
// 定义列标题
headers := []string{
    i18n.T("CONTAINER ID"),
    i18n.T("IMAGE"),
    i18n.T("COMMAND"),
    i18n.T("CREATED"),
    i18n.T("STATUS"),
}

// 计算列宽（考虑多字节字符）
columnWidths := make([]int, len(headers))
for i, header := range headers {
    columnWidths[i] = i18n.GetDisplayWidth(header) + 2  // +2 for padding
}
```

### 场景 4: 进度消息

```go
// 显示进度消息
progressMsg := i18n.T("adding new connection")
fmt.Println(progressMsg)

// 或者在 TUI 中
dialog.SetText(i18n.T("podman disk usage in progress"))
```

### 场景 5: 错误消息

```go
// 错误消息
if err != nil {
    return fmt.Errorf("%s: %w", i18n.T("failed to connect"), err)
}
```

### 场景 6: 帮助文本

```go
// 命令描述
commands := map[string]string{
    "connect":    i18n.T("connect to selected destination"),
    "disconnect": i18n.T("disconnect from connected destination"),
    "info":       i18n.T("display destination system information"),
}
```

## 代码迁移步骤

### 步骤 1: 识别需要翻译的字符串

查找所有面向用户的字符串：
- 按钮文本
- 菜单项
- 标签
- 错误消息
- 帮助文本
- 状态消息

### 步骤 2: 添加翻译条目

在 `i18n/zh_cn.go` 和其他语言文件中添加对应的翻译：

```go
var ZhCnTranslations = map[string]string{
    // 新增的翻译
    "Your New String": "你的新字符串",
    // ...
}
```

### 步骤 3: 替换代码中的字符串

使用 `i18n.T()` 包装所有需要翻译的字符串。

### 步骤 4: 修复宽度计算

将所有使用 `len()` 计算显示宽度的地方替换为 `i18n.GetDisplayWidth()`。

### 步骤 5: 测试

在不同语言环境下测试：

```bash
# 测试中文
LANG=zh_CN.UTF-8 ./bin/podman-tui

# 测试日文
LANG=ja_JP.UTF-8 ./bin/podman-tui

# 测试英文
LANG=en_US.UTF-8 ./bin/podman-tui
```

## 实际代码示例

### 示例 1: 修改对话框

**之前：**
```go
func showConfirmDialog(title, message string) bool {
    dialog := tview.NewModal().
        SetText(message).
        AddButtons([]string{"OK", "Cancel"})
    // ...
}
```

**之后：**
```go
func showConfirmDialog(title, message string) bool {
    dialog := tview.NewModal().
        SetText(i18n.T(message)).
        AddButtons([]string{i18n.T("OK"), i18n.T("Cancel")})
    // ...
}
```

### 示例 2: 修改表格

**之前：**
```go
table.SetCell(0, 0, tview.NewTableCell("Name"))
table.SetCell(0, 1, tview.NewTableCell("Status"))
table.SetCell(0, 2, tview.NewTableCell("Image"))
```

**之后：**
```go
table.SetCell(0, 0, tview.NewTableCell(i18n.T("Name")))
table.SetCell(0, 1, tview.NewTableCell(i18n.T("Status")))
table.SetCell(0, 2, tview.NewTableCell(i18n.T("Image")))
```

### 示例 3: 修改表单

**之前：**
```go
form := tview.NewForm().
    AddInputField("Container Name:", "", 20, nil, nil).
    AddInputField("Image:", "", 20, nil, nil).
    AddButton("Create", createHandler).
    AddButton("Cancel", cancelHandler)
```

**之后：**
```go
form := tview.NewForm().
    AddInputField(i18n.T("Container Name")+":", "", 20, nil, nil).
    AddInputField(i18n.T("Image")+":", "", 20, nil, nil).
    AddButton(i18n.T("Create"), createHandler).
    AddButton(i18n.T("Cancel"), cancelHandler)
```

### 示例 4: 修复宽度问题

**之前（有 bug）：**
```go
// 这在中文环境下会导致对齐问题
labelWidth := len("容器名称:") + 1  // 错误：返回 13
paddedLabel := fmt.Sprintf("%-*s", labelWidth, "容器名称:")
```

**之后（正确）：**
```go
label := i18n.T("Container Name") + ":"
labelWidth := i18n.GetDisplayWidth(label) + 1
paddedLabel := i18n.PadToWidth(label, labelWidth)
```

## 性能考虑

### 1. 缓存翻译结果

对于频繁使用的字符串，可以缓存翻译结果：

```go
type MyWidget struct {
    cancelText string
    okText     string
}

func NewMyWidget() *MyWidget {
    return &MyWidget{
        cancelText: i18n.T("Cancel"),
        okText:     i18n.T("OK"),
    }
}
```

### 2. 避免重复翻译

**不推荐：**
```go
for i := 0; i < 1000; i++ {
    label := i18n.T("Status")  // 每次循环都查找
    // ...
}
```

**推荐：**
```go
statusLabel := i18n.T("Status")  // 只查找一次
for i := 0; i < 1000; i++ {
    label := statusLabel
    // ...
}
```

## 调试技巧

### 1. 检查当前语言

```go
fmt.Printf("Current language: %s\n", i18n.GetCurrentLanguage())
```

### 2. 强制使用特定语言

```go
// 用于测试
i18n.SetLanguage("zh_cn")
```

### 3. 查找缺失的翻译

如果某个字符串没有翻译，`i18n.T()` 会返回原始的英文字符串，这样可以很容易发现缺失的翻译。

## 常见问题

### Q: 如何处理带参数的字符串？

A: 目前使用 `fmt.Sprintf`：

```go
msg := fmt.Sprintf("%s: %d", i18n.T("Total containers"), count)
```

### Q: 如何处理复数形式？

A: 目前需要手动处理：

```go
var msg string
if count == 1 {
    msg = i18n.T("1 container")
} else {
    msg = fmt.Sprintf("%d %s", count, i18n.T("containers"))
}
```

### Q: 翻译文件太大怎么办？

A: 可以考虑按模块分割翻译文件，例如：
- `zh_cn_containers.go`
- `zh_cn_images.go`
- `zh_cn_volumes.go`

### Q: 如何确保所有语言的翻译一致？

A: 可以编写一个工具来检查所有语言文件是否包含相同的键：

```bash
# 提取所有键
grep -o '"[^"]*":' i18n/zh_cn.go | sort > keys_zh_cn.txt
grep -o '"[^"]*":' i18n/ja.go | sort > keys_ja.txt

# 比较
diff keys_zh_cn.txt keys_ja.txt
```

## 下一步

1. **逐步迁移**: 建议从一个模块开始，逐步迁移到 i18n
2. **测试覆盖**: 确保在不同语言环境下测试 UI
3. **文档更新**: 更新用户文档，说明支持的语言
4. **贡献指南**: 鼓励社区贡献翻译

## 参考资源

- [i18n/README.md](README.md) - i18n 框架使用文档
- [i18n/example/main.go](example/main.go) - 完整示例程序
- [i18n/i18n_test.go](i18n_test.go) - 单元测试
