package main

import (
	"fmt"

	"github.com/containers/podman-tui/i18n"
)

func main() {
	fmt.Println("=== Podman-TUI i18n Framework Demo ===\n")

	// 1. 显示当前语言（默认为英文）
	currentLang := i18n.GetCurrentLanguage()
	fmt.Printf("Current language: %s (%s)\n", currentLang, i18n.GetLanguageName(currentLang))
	fmt.Printf("Note: Default is English. Use i18n.SetLanguage() to change.\n\n")

	// 2. 显示所有支持的语言
	fmt.Println("Supported languages:")
	for _, lang := range i18n.GetSupportedLanguages() {
		fmt.Printf("  - %s: %s\n", lang, i18n.GetLanguageName(lang))
	}
	fmt.Println()

	// 3. 演示翻译功能
	fmt.Println("Translation examples:")
	testKeys := []string{"Cancel", "HELP", "add connection", "delete", "info"}
	
	for _, key := range testKeys {
		fmt.Printf("  '%s' → '%s'\n", key, i18n.T(key))
	}
	fmt.Println()

	// 4. 演示切换语言
	fmt.Println("Switching languages:")
	languages := []string{"en", "zh_cn", "ja", "ko", "ru"}
	
	for _, lang := range languages {
		i18n.SetLanguage(lang)
		fmt.Printf("  [%s] Cancel = %s\n", i18n.GetLanguageName(lang), i18n.T("Cancel"))
	}
	fmt.Println()

	// 5. 演示字符宽度处理
	i18n.SetLanguage("zh_cn")
	fmt.Println("Display width examples (Chinese):")
	
	testStrings := []string{"删除", "Delete", "容器名称", "Container Name"}
	for _, str := range testStrings {
		width := i18n.GetDisplayWidth(str)
		fmt.Printf("  '%s' → width: %d (len: %d)\n", str, width, len(str))
	}
	fmt.Println()

	// 6. 演示字符串填充
	fmt.Println("Padding examples:")
	fmt.Printf("  English: '%s'\n", i18n.PadToWidth("Delete", 20))
	fmt.Printf("  Chinese: '%s'\n", i18n.PadToWidth("删除", 20))
	fmt.Println()

	// 7. 演示标签格式化
	fmt.Println("Label formatting examples:")
	labels := []string{"Name", "Status", "Image", "Created"}
	i18n.SetLanguage("en")
	fmt.Println("  English:")
	for _, label := range labels {
		formatted := i18n.FormatLabelWithWidth(label, 15, ":")
		fmt.Printf("    %s value\n", formatted)
	}
	
	i18n.SetLanguage("zh_cn")
	fmt.Println("  Chinese:")
	chineseLabels := []string{"名称", "状态", "镜像", "创建时间"}
	for _, label := range chineseLabels {
		formatted := i18n.FormatLabelWithWidth(label, 15, ":")
		fmt.Printf("    %s 值\n", formatted)
	}
	fmt.Println()

	// 8. 演示宽度调整
	fmt.Println("Width adjustment for locale:")
	baseWidth := 20
	i18n.SetLanguage("en")
	fmt.Printf("  English base width: %d → adjusted: %d\n", baseWidth, i18n.AdjustWidthForLocale(baseWidth))
	i18n.SetLanguage("zh_cn")
	fmt.Printf("  Chinese base width: %d → adjusted: %d\n", baseWidth, i18n.AdjustWidthForLocale(baseWidth))
	fmt.Println()

	// 9. 演示截断功能
	fmt.Println("Truncation examples:")
	longText := "这是一个非常非常长的中文字符串用于测试截断功能"
	fmt.Printf("  Original: '%s' (width: %d)\n", longText, i18n.GetDisplayWidth(longText))
	fmt.Printf("  Truncated to 20: '%s'\n", i18n.TruncateToWidth(longText, 20))
	fmt.Printf("  Truncated to 10: '%s'\n", i18n.TruncateToWidth(longText, 10))
}
