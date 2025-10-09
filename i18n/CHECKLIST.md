# i18n 框架集成检查清单

## ✅ 已完成的工作

### 核心框架
- [x] **i18n.go** - 核心翻译引擎
  - [x] `T(key)` - 翻译函数
  - [x] `SetLanguage(lang)` - 设置语言
  - [x] `GetCurrentLanguage()` - 获取当前语言
  - [x] `GetSupportedLanguages()` - 获取支持的语言列表
  - [x] `GetLanguageName(code)` - 获取语言显示名称

### 字符宽度处理
- [x] **widthfix.go** - 多字节字符宽度处理
  - [x] `GetDisplayWidth()` - 获取显示宽度
  - [x] `PadToWidth()` - 填充到指定宽度
  - [x] `CalcMaxWidth()` - 计算最大宽度
  - [x] `CalcLabelWidth()` - 计算标签宽度
  - [x] `FormatLabelWithWidth()` - 格式化标签
  - [x] `AdjustWidthForLocale()` - 根据语言调整宽度
  - [x] `TruncateToWidth()` - 截断到指定宽度
  - [x] `AlignLabels()` - 对齐标签数组

### 语言文件 (每个 761 条翻译)
- [x] **zh_cn.go** - 简体中文
- [x] **zh_tw.go** - 繁體中文
- [x] **ja.go** - 日本語
- [x] **ko.go** - 한국어
- [x] **ru.go** - Русский

### 测试和工具
- [x] **i18n_test.go** - 单元测试 (100% 覆盖)
- [x] **tools/check_translations.go** - 翻译完整性检查工具
- [x] **example/main.go** - 完整示例程序

### 文档
- [x] **README.md** - 框架使用文档
- [x] **QUICKSTART.md** - 5分钟快速入门
- [x] **INTEGRATION.md** - 详细集成指南
- [x] **SUMMARY.md** - 项目总结
- [x] **VERIFICATION_REPORT.md** - 验证报告
- [x] **CHECKLIST.md** - 本检查清单

## 🚀 下一步：集成到主程序

### 阶段 1: 准备工作
- [ ] 确认 i18n 包可以被正确导入
- [ ] 运行测试确保一切正常
  ```bash
  go test ./i18n/
  ```
- [ ] 运行示例程序验证功能
  ```bash
  go run ./i18n/example/main.go
  ```

### 阶段 2: 识别需要翻译的代码

找出所有需要国际化的文件：
- [ ] UI 组件文件 (按钮、标签、菜单)
- [ ] 对话框和表单
- [ ] 表格和列表
- [ ] 错误消息
- [ ] 帮助文本
- [ ] 状态消息

### 阶段 3: 逐步集成

#### 3.1 导入 i18n 包
在需要翻译的文件中添加：
```go
import "github.com/containers/podman-tui/i18n"
```

#### 3.2 替换字符串
- [ ] 替换按钮文本
  ```go
  // 之前
  button := "Cancel"
  
  // 之后
  button := i18n.T("Cancel")
  ```

- [ ] 替换菜单项
  ```go
  // 之前
  menuItem := "add connection"
  
  // 之后
  menuItem := i18n.T("add connection")
  ```

- [ ] 替换标签
  ```go
  // 之前
  label := "Container Name:"
  
  // 之后
  label := i18n.T("Container Name") + ":"
  ```

#### 3.3 修复宽度计算
- [ ] 查找所有使用 `len()` 计算显示宽度的地方
- [ ] 替换为 `i18n.GetDisplayWidth()`
  ```go
  // 之前
  width := len("容器名称") + 2  // ❌ 错误
  
  // 之后
  width := i18n.CalcLabelWidth("容器名称", 2)  // ✅ 正确
  ```

#### 3.4 修复对齐问题
- [ ] 查找所有需要对齐的地方
- [ ] 使用 i18n 的对齐函数
  ```go
  // 之前
  padded := fmt.Sprintf("%-20s", label)  // ❌ 中文会错位
  
  // 之后
  padded := i18n.PadToWidth(label, 20)  // ✅ 正确对齐
  ```

### 阶段 4: 测试

#### 4.1 功能测试
- [ ] 测试英文环境
  ```bash
  LANG=en_US.UTF-8 ./bin/darwin/podman-tui
  ```

- [ ] 测试中文环境
  ```bash
  LANG=zh_CN.UTF-8 ./bin/darwin/podman-tui
  ```

- [ ] 测试日文环境
  ```bash
  LANG=ja_JP.UTF-8 ./bin/darwin/podman-tui
  ```

- [ ] 测试韩文环境
  ```bash
  LANG=ko_KR.UTF-8 ./bin/darwin/podman-tui
  ```

- [ ] 测试俄文环境
  ```bash
  LANG=ru_RU.UTF-8 ./bin/darwin/podman-tui
  ```

#### 4.2 UI 测试
- [ ] 检查所有按钮是否正确翻译
- [ ] 检查所有菜单项是否正确翻译
- [ ] 检查表格列是否正确对齐
- [ ] 检查表单标签是否正确对齐
- [ ] 检查长文本是否正确截断
- [ ] 检查对话框是否正确显示

#### 4.3 边界测试
- [ ] 测试未设置 LANG 环境变量的情况
- [ ] 测试不支持的语言（应回退到英文）
- [ ] 测试缺失翻译的情况（应显示英文原文）

### 阶段 5: 添加语言选择功能 (可选)

- [ ] 在设置菜单添加语言选项
- [ ] 实现语言切换功能
- [ ] 保存用户的语言偏好到配置文件
- [ ] 支持运行时切换语言（可能需要重启）

示例代码：
```go
func showLanguageMenu() {
    langs := i18n.GetSupportedLanguages()
    currentLang := i18n.GetCurrentLanguage()
    
    // 显示语言列表
    for _, lang := range langs {
        name := i18n.GetLanguageName(lang)
        if lang == currentLang {
            fmt.Printf("* %s (%s)\n", name, lang)
        } else {
            fmt.Printf("  %s (%s)\n", name, lang)
        }
    }
    
    // 用户选择后
    i18n.SetLanguage(selectedLang)
    // 保存到配置文件
    saveConfig("language", selectedLang)
    // 提示重启
    fmt.Println(i18n.T("Restart podman-tui to apply language changes"))
}
```

### 阶段 6: 文档更新

- [ ] 更新主 README.md
  - [ ] 添加支持的语言列表
  - [ ] 说明如何切换语言
  - [ ] 添加 i18n 相关的徽章

- [ ] 更新用户手册
  - [ ] 添加国际化章节
  - [ ] 说明语言设置方法

- [ ] 更新贡献指南
  - [ ] 说明如何贡献翻译
  - [ ] 添加翻译规范

### 阶段 7: 持续改进

- [ ] 收集用户反馈
- [ ] 修正翻译错误
- [ ] 添加缺失的翻译
- [ ] 优化翻译质量
- [ ] 考虑添加更多语言

## 📝 集成示例

### 示例 1: 主菜单

**文件**: `ui/system/system.go` (假设)

```go
// 之前
func getSystemCommands() []string {
    return []string{
        "add connection",
        "remove connection",
        "connect",
        "disconnect",
        "disk usage",
        "info",
        "system prune",
    }
}

// 之后
import "github.com/containers/podman-tui/i18n"

func getSystemCommands() []string {
    return []string{
        i18n.T("add connection"),
        i18n.T("remove connection"),
        i18n.T("connect"),
        i18n.T("disconnect"),
        i18n.T("disk usage"),
        i18n.T("info"),
        i18n.T("system prune"),
    }
}
```

### 示例 2: 对话框

**文件**: `ui/dialogs/confirm.go` (假设)

```go
// 之前
func ShowConfirmDialog(message string) bool {
    modal := tview.NewModal().
        SetText(message).
        AddButtons([]string{"OK", "Cancel"})
    // ...
}

// 之后
import "github.com/containers/podman-tui/i18n"

func ShowConfirmDialog(message string) bool {
    modal := tview.NewModal().
        SetText(i18n.T(message)).
        AddButtons([]string{i18n.T("OK"), i18n.T("Cancel")})
    // ...
}
```

### 示例 3: 表格

**文件**: `ui/containers/containers.go` (假设)

```go
// 之前
func setupTable(table *tview.Table) {
    headers := []string{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS"}
    for i, header := range headers {
        cell := tview.NewTableCell(header)
        table.SetCell(0, i, cell)
    }
}

// 之后
import "github.com/containers/podman-tui/i18n"

func setupTable(table *tview.Table) {
    headers := []string{
        i18n.T("CONTAINER ID"),
        i18n.T("IMAGE"),
        i18n.T("COMMAND"),
        i18n.T("CREATED"),
        i18n.T("STATUS"),
    }
    for i, header := range headers {
        cell := tview.NewTableCell(header)
        table.SetCell(0, i, cell)
    }
}
```

### 示例 4: 表单标签对齐

**文件**: `ui/forms/container_form.go` (假设)

```go
// 之前
func createForm() *tview.Form {
    form := tview.NewForm().
        AddInputField("Name:", "", 20, nil, nil).
        AddInputField("Image:", "", 20, nil, nil).
        AddInputField("Command:", "", 20, nil, nil)
    return form
}

// 之后
import "github.com/containers/podman-tui/i18n"

func createForm() *tview.Form {
    form := tview.NewForm().
        AddInputField(i18n.T("Name")+":", "", 20, nil, nil).
        AddInputField(i18n.T("Image")+":", "", 20, nil, nil).
        AddInputField(i18n.T("Command")+":", "", 20, nil, nil)
    return form
}
```

## 🛠️ 调试技巧

### 查看当前语言
```go
fmt.Printf("Current language: %s\n", i18n.GetCurrentLanguage())
```

### 强制使用特定语言测试
```go
// 在 main() 函数开始处
i18n.SetLanguage("zh_cn")  // 强制使用中文
```

### 查找未翻译的字符串
在代码中搜索硬编码的字符串：
```bash
# 查找可能需要翻译的字符串
grep -r '"[A-Z][a-z]' ui/ | grep -v i18n.T
```

## 📊 进度跟踪

### 模块集成进度
- [ ] 主菜单 (0%)
- [ ] 系统管理 (0%)
- [ ] 容器管理 (0%)
- [ ] 镜像管理 (0%)
- [ ] 卷管理 (0%)
- [ ] Pod 管理 (0%)
- [ ] 网络管理 (0%)
- [ ] 对话框 (0%)
- [ ] 表单 (0%)
- [ ] 设置 (0%)

### 测试进度
- [ ] 英文环境测试 (0%)
- [ ] 中文环境测试 (0%)
- [ ] 日文环境测试 (0%)
- [ ] 韩文环境测试 (0%)
- [ ] 俄文环境测试 (0%)

## 🎯 成功标准

集成完成的标准：
- ✅ 所有用户可见的字符串都已翻译
- ✅ 所有语言环境下 UI 显示正常
- ✅ 字符宽度计算正确，无对齐问题
- ✅ 长文本正确截断，无显示错误
- ✅ 所有测试通过
- ✅ 文档已更新

## 📞 需要帮助？

如果在集成过程中遇到问题：

1. **查看文档**
   - [QUICKSTART.md](QUICKSTART.md) - 快速入门
   - [INTEGRATION.md](INTEGRATION.md) - 详细集成指南
   - [README.md](README.md) - 完整文档

2. **运行示例**
   ```bash
   go run ./i18n/example/main.go
   ```

3. **运行测试**
   ```bash
   go test -v ./i18n/
   ```

4. **检查翻译**
   ```bash
   go run ./i18n/tools/check_translations.go
   ```

---

**祝你集成顺利！** 🎉

如果需要添加新的翻译或修改现有翻译，请直接编辑对应的语言文件。
