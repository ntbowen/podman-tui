# i18n 工具

这个目录包含用于管理和验证翻译的工具。

## check_translations.go

检查所有语言文件的翻译完整性和一致性。

### 使用方法

```bash
# 从 i18n 目录运行
cd /path/to/podman-tui/i18n
go run tools/check_translations.go

# 或者指定 i18n 目录路径
go run tools/check_translations.go /path/to/i18n/
```

### 功能

1. **统计每个语言文件的翻译条目数量**
2. **找出参考语言**（键最多的语言）
3. **比较所有语言与参考语言**
   - 列出缺失的翻译键
   - 列出多余的翻译键
4. **报告不一致的地方**

### 示例输出

```
=== Translation Completeness Check ===

✓ zh_cn: 150 keys
✓ zh_tw: 150 keys
✓ ja: 148 keys
✓ ko: 145 keys
✓ ru: 150 keys

Using zh_cn as reference (150 keys)
==================================================

## Comparing zh_tw with zh_cn

✓ zh_tw has the same number of keys as zh_cn

## Comparing ja with zh_cn

❌ Keys in zh_cn but missing in ja (2):
   - "new feature 1"
   - "new feature 2"

⚠️  ja is missing 2 translations

## Comparing ko with zh_cn

❌ Keys in zh_cn but missing in ko (5):
   - "feature A"
   - "feature B"
   - "feature C"
   - "feature D"
   - "feature E"

⚠️  ko is missing 5 translations

==================================================

⚠️  Some translations are incomplete or inconsistent.
Please review the differences above.
```

## 未来工具计划

### 1. extract_strings.go
从源代码中提取所有 `i18n.T()` 调用，生成需要翻译的字符串列表。

### 2. generate_template.go
生成新语言的翻译模板文件。

### 3. sync_translations.go
同步所有语言文件，确保它们包含相同的键。

### 4. validate_format.go
验证翻译文件的格式是否正确。

### 5. coverage_report.go
生成翻译覆盖率报告，显示哪些源文件使用了 i18n，哪些还没有。
