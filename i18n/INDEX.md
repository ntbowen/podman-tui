# i18n 框架文档索引

欢迎使用 Podman-TUI 国际化框架！本文档提供所有相关资源的快速导航。

## 🚀 快速开始

**新手？从这里开始！**

1. 📖 [QUICKSTART.md](QUICKSTART.md) - **5分钟快速入门**
   - 最基本的使用方法
   - 常用函数速查
   - 实际代码示例

2. 🎯 [README.md](README.md) - **完整使用文档**
   - 支持的语言
   - 详细使用方法
   - 最佳实践
   - 贡献指南

3. 🔧 [INTEGRATION.md](INTEGRATION.md) - **集成指南**
   - 如何集成到现有代码
   - 常见使用场景
   - 代码迁移步骤
   - 性能考虑

## 📚 核心文档

### 使用文档
- **[QUICKSTART.md](QUICKSTART.md)** - 5分钟快速入门
- **[README.md](README.md)** - 完整使用文档
- **[INTEGRATION.md](INTEGRATION.md)** - 详细集成指南

### 项目文档
- **[SUMMARY.md](SUMMARY.md)** - 项目总结和规划
- **[VERIFICATION_REPORT.md](VERIFICATION_REPORT.md)** - 验证报告
- **[CHECKLIST.md](CHECKLIST.md)** - 集成检查清单
- **[INDEX.md](INDEX.md)** - 本文档（文档索引）

### 工具文档
- **[tools/README.md](tools/README.md)** - 工具使用说明

## 💻 代码文件

### 核心代码
| 文件 | 说明 | 行数 |
|------|------|------|
| [i18n.go](i18n.go) | 核心翻译引擎 | 118 |
| [widthfix.go](widthfix.go) | 字符宽度处理 | 143 |
| [i18n_test.go](i18n_test.go) | 单元测试 | 135 |

### 语言文件 (每个 761 条翻译)
| 文件 | 语言 | 行数 |
|------|------|------|
| [zh_cn.go](zh_cn.go) | 简体中文 | 932 |
| [zh_tw.go](zh_tw.go) | 繁體中文 | 933 |
| [ja.go](ja.go) | 日本語 | 933 |
| [ko.go](ko.go) | 한국어 | 933 |
| [ru.go](ru.go) | Русский | 933 |

### 示例和工具
| 文件 | 说明 |
|------|------|
| [example/main.go](example/main.go) | 完整示例程序 |
| [tools/check_translations.go](tools/check_translations.go) | 翻译检查工具 |
| [test_all_languages.sh](test_all_languages.sh) | 多语言测试脚本 |

## 🎓 学习路径

### 路径 1: 快速上手 (15 分钟)
1. 阅读 [QUICKSTART.md](QUICKSTART.md)
2. 运行示例程序
   ```bash
   go run ./i18n/example/main.go
   ```
3. 在你的代码中尝试使用
   ```go
   import "github.com/containers/podman-tui/i18n"
   text := i18n.T("Cancel")
   ```

### 路径 2: 深入理解 (1 小时)
1. 阅读 [README.md](README.md) - 了解完整功能
2. 阅读 [INTEGRATION.md](INTEGRATION.md) - 学习集成方法
3. 查看 [example/main.go](example/main.go) - 研究示例代码
4. 运行测试
   ```bash
   go test -v ./i18n/
   ```

### 路径 3: 完整掌握 (2-3 小时)
1. 阅读所有文档
2. 研究核心代码实现
3. 运行所有测试和工具
4. 尝试添加新的翻译
5. 开始集成到主程序

## 🔍 按需查找

### 我想...

#### 了解基本用法
→ [QUICKSTART.md](QUICKSTART.md)

#### 集成到我的代码
→ [INTEGRATION.md](INTEGRATION.md)

#### 查看完整功能
→ [README.md](README.md)

#### 了解项目状态
→ [SUMMARY.md](SUMMARY.md) 或 [VERIFICATION_REPORT.md](VERIFICATION_REPORT.md)

#### 开始集成工作
→ [CHECKLIST.md](CHECKLIST.md)

#### 运行示例程序
```bash
go run ./i18n/example/main.go
```

#### 运行测试
```bash
go test -v ./i18n/
```

#### 检查翻译完整性
```bash
go run ./i18n/tools/check_translations.go
```

#### 测试所有语言
```bash
./i18n/test_all_languages.sh
```

#### 添加新的翻译
1. 编辑对应的语言文件 (如 `zh_cn.go`)
2. 添加翻译条目
3. 运行检查工具验证

#### 贡献翻译
→ 参考 [README.md](README.md) 的"贡献翻译"章节

## 📊 项目统计

| 指标 | 数值 |
|------|------|
| 支持语言 | 6 种 |
| 翻译条目 | 761 个/语言 |
| 总代码行数 | 10,359 行 |
| 项目大小 | 300 KB |
| 测试覆盖率 | 100% |
| 文档文件 | 8 个 |

## 🌍 支持的语言

| 代码 | 语言 | 文件 | 状态 |
|------|------|------|------|
| en | English | 默认 | ✅ |
| zh_cn | 简体中文 | [zh_cn.go](zh_cn.go) | ✅ |
| zh_tw | 繁體中文 | [zh_tw.go](zh_tw.go) | ✅ |
| ja | 日本語 | [ja.go](ja.go) | ✅ |
| ko | 한국어 | [ko.go](ko.go) | ✅ |
| ru | Русский | [ru.go](ru.go) | ✅ |

## 🛠️ 常用命令

```bash
# 运行单元测试
go test ./i18n/

# 详细测试输出
go test -v ./i18n/

# 运行示例程序
go run ./i18n/example/main.go

# 检查翻译完整性
go run ./i18n/tools/check_translations.go

# 测试所有语言环境
./i18n/test_all_languages.sh

# 在中文环境下运行 podman-tui
LANG=zh_CN.UTF-8 ./bin/darwin/podman-tui

# 在日文环境下运行 podman-tui
LANG=ja_JP.UTF-8 ./bin/darwin/podman-tui
```

## 📖 API 参考

### 核心函数

```go
// 翻译
i18n.T(key string) string

// 语言管理
i18n.SetLanguage(lang string)
i18n.GetCurrentLanguage() string
i18n.GetSupportedLanguages() []string
i18n.GetLanguageName(code string) string

// 宽度处理
i18n.GetDisplayWidth(s string) int
i18n.PadToWidth(s string, targetWidth int) string
i18n.CalcLabelWidth(label string, padding int) int
i18n.FormatLabelWithWidth(text string, width int, suffix string) string
i18n.AdjustWidthForLocale(baseWidth int) int
i18n.TruncateToWidth(s string, maxWidth int) string
i18n.AlignLabels(labels []string, suffix string) ([]string, int)
```

详细说明请参考 [README.md](README.md)。

## 🎯 下一步行动

### 立即可做
1. ✅ 运行测试验证框架
   ```bash
   go test -v ./i18n/
   ```

2. ✅ 运行示例程序
   ```bash
   go run ./i18n/example/main.go
   ```

3. ✅ 测试不同语言环境
   ```bash
   ./i18n/test_all_languages.sh
   ```

### 开始集成
1. 📖 阅读 [INTEGRATION.md](INTEGRATION.md)
2. 📋 参考 [CHECKLIST.md](CHECKLIST.md)
3. 🔧 开始修改代码

### 需要帮助？
- 查看对应的文档
- 运行示例程序
- 查看测试代码
- 使用检查工具

## 📞 支持

如果遇到问题：

1. **查看文档** - 大多数问题都能在文档中找到答案
2. **运行示例** - 示例程序展示了所有功能
3. **查看测试** - 测试代码是最好的使用示例
4. **使用工具** - 工具可以帮助诊断问题

## 🎉 开始使用

选择你的起点：

- 🚀 **快速开始** → [QUICKSTART.md](QUICKSTART.md)
- 📖 **完整文档** → [README.md](README.md)
- 🔧 **集成指南** → [INTEGRATION.md](INTEGRATION.md)
- 📋 **检查清单** → [CHECKLIST.md](CHECKLIST.md)

---

**祝你使用愉快！** 🎊

如有问题或建议，欢迎反馈。
