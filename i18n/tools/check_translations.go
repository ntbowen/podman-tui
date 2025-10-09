package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// 从 Go 文件中提取翻译键
func extractKeys(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 匹配 "key": "value" 格式
	re := regexp.MustCompile(`"([^"]+)":\s*"[^"]*"`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	keys := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			keys = append(keys, match[1])
		}
	}

	return keys, nil
}

// 比较两个键列表，返回差异
func compareKeys(keys1, keys2 []string, name1, name2 string) {
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, key := range keys1 {
		set1[key] = true
	}
	for _, key := range keys2 {
		set2[key] = true
	}

	// 找出在 set1 中但不在 set2 中的键
	missing := []string{}
	for key := range set1 {
		if !set2[key] {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		fmt.Printf("\n❌ Keys in %s but missing in %s (%d):\n", name1, name2, len(missing))
		for _, key := range missing {
			fmt.Printf("   - \"%s\"\n", key)
		}
	}
}

func main() {
	i18nDir := "../"
	if len(os.Args) > 1 {
		i18nDir = os.Args[1]
	}

	// 语言文件列表
	langFiles := map[string]string{
		"en":    "en.go",      // 如果有的话
		"zh_cn": "zh_cn.go",
		"zh_tw": "zh_tw.go",
		"ja":    "ja.go",
		"ko":    "ko.go",
		"ru":    "ru.go",
	}

	fmt.Println("=== Translation Completeness Check ===\n")

	// 提取所有语言的键
	allKeys := make(map[string][]string)
	for lang, file := range langFiles {
		path := filepath.Join(i18nDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("⚠️  %s not found, skipping...\n", file)
			continue
		}

		keys, err := extractKeys(path)
		if err != nil {
			fmt.Printf("❌ Error reading %s: %v\n", file, err)
			continue
		}

		allKeys[lang] = keys
		fmt.Printf("✓ %s: %d keys\n", lang, len(keys))
	}

	if len(allKeys) == 0 {
		fmt.Println("\n❌ No translation files found!")
		os.Exit(1)
	}

	// 找出键最多的语言作为参考
	var refLang string
	var refKeys []string
	maxKeys := 0

	for lang, keys := range allKeys {
		if len(keys) > maxKeys {
			maxKeys = len(keys)
			refLang = lang
			refKeys = keys
		}
	}

	fmt.Printf("\nUsing %s as reference (%d keys)\n", refLang, len(refKeys))
	fmt.Println(strings.Repeat("=", 50))

	// 比较所有语言与参考语言
	hasIssues := false
	for lang, keys := range allKeys {
		if lang == refLang {
			continue
		}

		fmt.Printf("\n## Comparing %s with %s\n", lang, refLang)
		
		// 检查缺失的键
		compareKeys(refKeys, keys, refLang, lang)
		
		// 检查多余的键
		compareKeys(keys, refKeys, lang, refLang)

		if len(refKeys) != len(keys) {
			hasIssues = true
			diff := len(refKeys) - len(keys)
			if diff > 0 {
				fmt.Printf("\n⚠️  %s is missing %d translations\n", lang, diff)
			} else {
				fmt.Printf("\n⚠️  %s has %d extra translations\n", lang, -diff)
			}
		} else {
			fmt.Printf("\n✓ %s has the same number of keys as %s\n", lang, refLang)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	
	if !hasIssues {
		fmt.Println("\n✅ All translations are complete and consistent!")
		os.Exit(0)
	} else {
		fmt.Println("\n⚠️  Some translations are incomplete or inconsistent.")
		fmt.Println("Please review the differences above.")
		os.Exit(1)
	}
}
