package i18n

import (
	"testing"
)

func TestT(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		key      string
		expected string
	}{
		{"English default", "en", "Cancel", "Cancel"},
		{"Chinese translation", "zh_cn", "Cancel", "取消"},
		{"Japanese translation", "ja", "Cancel", "キャンセル"},
		{"Korean translation", "ko", "Cancel", "취소"},
		{"Russian translation", "ru", "Cancel", "Отмена"},
		{"Traditional Chinese", "zh_tw", "Cancel", "取消"},
		{"Missing key fallback", "zh_cn", "NonExistentKey", "NonExistentKey"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original language
			originalLang := currentLanguage
			defer func() { currentLanguage = originalLang }()

			// Set test language
			SetLanguage(tt.lang)

			result := T(tt.key)
			if result != tt.expected {
				t.Errorf("T(%q) with lang %q = %v, want %v", tt.key, tt.lang, result, tt.expected)
			}
		})
	}
}

func TestGetSupportedLanguages(t *testing.T) {
	langs := GetSupportedLanguages()
	expectedCount := 6 // en, zh_cn, zh_tw, ja, ko, ru

	if len(langs) != expectedCount {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(langs), expectedCount)
	}

	// Check if all expected languages are present
	expectedLangs := map[string]bool{
		"en": true, "zh_cn": true, "zh_tw": true,
		"ja": true, "ko": true, "ru": true,
	}

	for _, lang := range langs {
		if !expectedLangs[lang] {
			t.Errorf("Unexpected language in supported list: %s", lang)
		}
	}
}

func TestGetLanguageName(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{"en", "English"},
		{"zh_cn", "简体中文"},
		{"zh_tw", "繁體中文"},
		{"ja", "日本語"},
		{"ko", "한국어"},
		{"ru", "Русский"},
		{"unknown", "unknown"}, // Fallback to code itself
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := GetLanguageName(tt.code)
			if result != tt.expected {
				t.Errorf("GetLanguageName(%q) = %v, want %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestSetLanguage(t *testing.T) {
	// Save original language
	originalLang := currentLanguage
	defer func() { currentLanguage = originalLang }()

	testLang := "zh_cn"
	SetLanguage(testLang)

	if currentLanguage != testLang {
		t.Errorf("SetLanguage(%q) failed, currentLanguage = %v", testLang, currentLanguage)
	}

	if GetCurrentLanguage() != testLang {
		t.Errorf("GetCurrentLanguage() = %v, want %v", GetCurrentLanguage(), testLang)
	}
}
