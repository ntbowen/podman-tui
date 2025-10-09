package i18n

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// GetDisplayWidth 返回字符串的实际显示宽度（支持中英文）
// 中文字符通常占2个显示位，英文字符占1个显示位
func GetDisplayWidth(s string) int {
	return runewidth.StringWidth(s)
}

// PadToWidth 将字符串填充到指定的显示宽度
// 如果字符串已经达到或超过目标宽度，则不添加填充
func PadToWidth(s string, targetWidth int) string {
	currentWidth := runewidth.StringWidth(s)
	if currentWidth >= targetWidth {
		return s
	}
	// 添加空格填充到目标宽度
	padding := targetWidth - currentWidth
	return s + strings.Repeat(" ", padding)
}

// CalcMaxWidth 计算多个字符串中的最大显示宽度
// 用于动态计算一组标签的最大宽度
func CalcMaxWidth(labels ...string) int {
	maxWidth := 0
	for _, label := range labels {
		width := runewidth.StringWidth(label)
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

// CalcLabelWidth 计算标签的显示宽度并添加填充
// 这是最常用的函数，用于替代 len() + 1 的用法
func CalcLabelWidth(label string, padding int) int {
	return runewidth.StringWidth(label) + padding
}

// FormatLabelWithWidth 格式化标签到指定宽度（右对齐）
// 例如：FormatLabelWithWidth("删除", 10, ":") -> "      删除:"
func FormatLabelWithWidth(text string, width int, suffix string) string {
	fullText := text + suffix
	currentWidth := runewidth.StringWidth(fullText)
	if currentWidth >= width {
		return fullText
	}
	// 计算需要添加的前导空格数
	padding := width - currentWidth
	return strings.Repeat(" ", padding) + fullText
}

// IsMultiByteLocale 判断当前是否为多字节字符环境（如中文）
// 通过检查当前语言是否为英文来判断
func IsMultiByteLocale() bool {
	return currentLanguage != "en_US"
}

// GetLabelWidthMultiplier 获取宽度倍数
// 对于中文等多字节语言，返回较大的倍数以确保足够的显示空间
func GetLabelWidthMultiplier() float64 {
	if IsMultiByteLocale() {
		// 中文字符平均占用约1.5倍的空间（考虑到混合使用中英文）
		return 1.5
	}
	return 1.0
}

// AdjustWidthForLocale 根据语言环境自动调整宽度
// baseWidth 是英文环境下的基准宽度
func AdjustWidthForLocale(baseWidth int) int {
	multiplier := GetLabelWidthMultiplier()
	adjustedWidth := int(float64(baseWidth) * multiplier)
	return adjustedWidth
}

// TruncateToWidth 将字符串截断到指定的显示宽度
// 如果字符串超过目标宽度，会添加省略号 "..."
func TruncateToWidth(s string, maxWidth int) string {
	currentWidth := runewidth.StringWidth(s)
	if currentWidth <= maxWidth {
		return s
	}
	
	// 如果需要截断，保留空间给省略号
	ellipsis := "..."
	ellipsisWidth := runewidth.StringWidth(ellipsis)
	targetWidth := maxWidth - ellipsisWidth
	
	if targetWidth <= 0 {
		return ellipsis
	}
	
	// 逐字符累加直到达到目标宽度
	result := []rune{}
	currentWidth = 0
	for _, r := range s {
		charWidth := runewidth.RuneWidth(r)
		if currentWidth+charWidth > targetWidth {
			break
		}
		result = append(result, r)
		currentWidth += charWidth
	}
	
	return string(result) + ellipsis
}

// AlignLabels 对齐一组标签到相同的宽度
// 返回对齐后的标签数组和使用的最大宽度
func AlignLabels(labels []string, suffix string) ([]string, int) {
	// 计算最大宽度
	maxWidth := 0
	for _, label := range labels {
		fullText := label + suffix
		width := runewidth.StringWidth(fullText)
		if width > maxWidth {
			maxWidth = width
		}
	}
	
	// 对齐所有标签
	aligned := make([]string, len(labels))
	for i, label := range labels {
		aligned[i] = PadToWidth(label+suffix, maxWidth)
	}
	
	return aligned, maxWidth
}

// GetSafeWidth 获取安全的显示宽度
// 在原有宽度基础上添加安全余量，避免被截断
func GetSafeWidth(s string, extraPadding int) int {
	baseWidth := runewidth.StringWidth(s)
	return baseWidth + extraPadding
}
