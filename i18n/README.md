# Podman-TUI 国际化 (i18n) 框架

这个目录包含了 podman-tui 的国际化支持框架。

## 支持的语言

- **en** - English (英语) - 默认语言
- **zh_cn** - 简体中文
- **zh_tw** - 繁體中文
- **ja** - 日本語 (日语)
- **ko** - 한국어 (韩语)
- **ru** - Русский (俄语)

## 文件结构

```
i18n/
├── i18n.go          # 核心 i18n 框架
├── widthfix.go      # 字符宽度处理工具(支持中文等多字节字符)
├── zh_cn.go         # 简体中文翻译
├── zh_tw.go         # 繁体中文翻译
├── ja.go            # 日语翻译
├── ko.go            # 韩语翻译
├── ru.go            # 俄语翻译
├── i18n_test.go     # 单元测试
└── README.md        # 本文档
```

## 使用方法

### 1. 在代码中使用翻译

```go
import "github.com/containers/podman-tui/i18n"

// 翻译字符串（默认英文，需手动设置语言）
text := i18n.T("Cancel")  // 返回 "Cancel"（默认）或 "取消"（设置中文后）

// 获取当前语言
lang := i18n.GetCurrentLanguage()  // 例如: "zh_cn"

// 手动设置语言
i18n.SetLanguage("ja")

// 获取支持的语言列表
langs := i18n.GetSupportedLanguages()  // ["en", "zh_cn", "zh_tw", "ja", "ko", "ru"]

// 获取语言的显示名称
name := i18n.GetLanguageName("zh_cn")  // "简体中文"
```

### 2. 处理多字节字符的显示宽度

由于中文、日文、韩文等字符在终端中占用2个显示位，而英文字符只占1个显示位，我们提供了专门的工具函数：

```go
import "github.com/containers/podman-tui/i18n"

// 获取字符串的实际显示宽度
width := i18n.GetDisplayWidth("删除")  // 返回 4 (2个中文字符)
width := i18n.GetDisplayWidth("Delete")  // 返回 6 (6个英文字符)

// 填充字符串到指定宽度
padded := i18n.PadToWidth("删除", 10)  // "删除      " (总宽度为10)

// 计算标签宽度(常用于表单)
labelWidth := i18n.CalcLabelWidth("容器名称", 2)  // 宽度 + 2

// 格式化标签(右对齐)
label := i18n.FormatLabelWithWidth("删除", 10, ":")  // "      删除:"

// 根据语言环境自动调整宽度
adjustedWidth := i18n.AdjustWidthForLocale(20)  // 中文环境返回 30

// 截断字符串到指定宽度
truncated := i18n.TruncateToWidth("这是一个很长的字符串", 10)  // "这是一..."
```

### 3. 添加新的翻译

在对应的语言文件中添加翻译条目：

```go
// zh_cn.go
var ZhCnTranslations = map[string]string{
    "New Key": "新的键",
    // ... 其他翻译
}
```

## 语言设置

框架默认使用英语。需要通过 `SetLanguage()` 手动设置语言：

```go
// 设置为中文
i18n.SetLanguage("zh_cn")

// 设置为日语
i18n.SetLanguage("ja")

// 获取当前语言
currentLang := i18n.GetCurrentLanguage()  // 返回 "en"（默认）
```

建议在应用启动时根据用户配置或系统环境设置语言。

## 测试

运行单元测试：

```bash
cd i18n
go test -v
```

## 最佳实践

### 1. 翻译键的命名

- 使用英文原文作为键，这样即使翻译缺失也能显示有意义的文本
- 保持键的简洁和一致性
- 对于长文本，可以使用简短的描述性键

```go
// 推荐
i18n.T("Cancel")
i18n.T("add connection")

// 不推荐
i18n.T("btn_cancel")
i18n.T("ADD_CONN")
```

### 2. 处理字符宽度

在需要对齐或计算显示宽度时，**务必使用** `widthfix.go` 中的函数，而不是 `len()`：

```go
// ❌ 错误 - len() 返回字节数，不是显示宽度
width := len("删除") + 1  // 返回 7，但实际显示宽度是 5

// ✅ 正确 - 使用 GetDisplayWidth
width := i18n.CalcLabelWidth("删除", 1)  // 返回 5
```

### 3. 翻译覆盖

- 确保所有语言文件包含相同的键
- 如果某个语言缺少翻译，系统会自动回退到英文原文
- 定期检查翻译的完整性

### 4. 上下文相关的翻译

对于同一个词在不同上下文中有不同翻译的情况，使用更具体的键：

```go
// 动词 "删除"
i18n.T("delete")

// 名词 "删除操作"
i18n.T("deletion")
```

## 贡献翻译

欢迎贡献新的语言或改进现有翻译！

1. 复制 `zh_cn.go` 作为模板
2. 重命名为新语言代码（如 `fr.go` 表示法语）
3. 翻译所有条目
4. 在 `i18n.go` 的 `T()` 函数中添加新语言的 case
5. 更新 `GetSupportedLanguages()` 和 `GetLanguageName()`
6. 提交 Pull Request

## 技术细节

### 依赖

- `github.com/mattn/go-runewidth` - 用于正确计算 Unicode 字符的显示宽度

### 设计原则

1. **简单明了** - 默认英文，手动设置语言
2. **优雅降级** - 翻译缺失时回退到英文原文
3. **类型安全** - 使用 map 存储翻译，编译时检查
4. **性能优化** - 翻译查找是 O(1) 的哈希表操作
5. **多字节字符支持** - 专门处理 CJK 字符的显示宽度问题

## 未来改进

- [ ] 支持复数形式 (pluralization)
- [ ] 支持参数化翻译 (如 "删除了 %d 个容器")
- [ ] 从外部文件加载翻译 (JSON/YAML)
- [ ] 翻译覆盖率检查工具
- [ ] 自动翻译建议 (通过 AI)
