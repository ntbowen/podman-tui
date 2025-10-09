package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const (
	// Language config file path
	_langConfigPath = "podman-tui/language.json"
)

// Current language (default to English)
var currentLanguage = "en"

// LanguageConfig stores the user's language preference
type LanguageConfig struct {
	Language string `json:"language"`
}

// T returns the translation for the given key
func T(key string) string {
	var translations map[string]string
	
	switch currentLanguage {
	case "zh_cn":
		translations = ZhCnTranslations
	case "zh_tw":
		translations = ZhTwTranslations
	case "ja":
		translations = JaTranslations
	case "ko":
		translations = KoTranslations
	case "ru":
		translations = RuTranslations
	case "en":
		// English is the default, return original text
		return key
	default:
		// Default to original text (English)
		return key
	}
	
	// Look up translation
	if translation, exists := translations[key]; exists {
		return translation
	}
	
	// Fallback to original text if translation not found
	return key
}

// GetCurrentLanguage returns the current language
func GetCurrentLanguage() string {
	return currentLanguage
}

// SetLanguage manually sets the language
func SetLanguage(lang string) {
	currentLanguage = lang
}

// GetSupportedLanguages returns a list of supported language codes
func GetSupportedLanguages() []string {
	return []string{"en", "zh_cn", "zh_tw", "ja", "ko", "ru"}
}

// GetLanguageName returns the display name for a language code
func GetLanguageName(code string) string {
	names := map[string]string{
		"en":    "English",
		"zh_cn": "简体中文",
		"zh_tw": "繁體中文",
		"ja":    "日本語",
		"ko":    "한국어",
		"ru":    "Русский",
	}
	
	if name, exists := names[code]; exists {
		return name
	}
	return code
}

// GetLanguageNameWithMark returns the display name with current language mark
func GetLanguageNameWithMark(code string) string {
	name := GetLanguageName(code)
	if code == currentLanguage {
		return "✅ " + name
	}
	return "   " + name
}

// LoadLanguageConfig loads the saved language preference
func LoadLanguageConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Warn().Msgf("failed to get user config dir: %v", err)
		return err
	}

	configPath := filepath.Join(configDir, _langConfigPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config file doesn't exist, use default language
			log.Debug().Msg("language config file not found, using default")
			return nil
		}
		return err
	}

	var config LanguageConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Warn().Msgf("failed to parse language config: %v", err)
		return err
	}

	if config.Language != "" {
		currentLanguage = config.Language
		log.Info().Msgf("loaded language preference: %s", currentLanguage)
	}

	return nil
}

// SaveLanguageConfig saves the current language preference
func SaveLanguageConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, _langConfigPath)
	configDirPath := filepath.Dir(configPath)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDirPath, 0o750); err != nil {
		return err
	}

	config := LanguageConfig{
		Language: currentLanguage,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return err
	}

	log.Info().Msgf("saved language preference: %s to %s", currentLanguage, configPath)
	return nil
}

// TranslateTime translates time duration strings (e.g., "5 seconds", "2 minutes", "23s", "42m38s")
func TranslateTime(duration string) string {
	if currentLanguage == "en" {
		return duration
	}
	
	// Short form translations (e.g., "23s" -> "23秒", "42m38s" -> "42分38秒")
	shortForms := map[string]string{
		"s": T("seconds"),
		"m": T("minutes"),
		"h": T("hours"),
		"d": T("days"),
		"w": T("weeks"),
	}
	
	// Time unit translations for full words
	// IMPORTANT: Must be ordered from longest to shortest to avoid partial matches
	// e.g., "minutes" must be checked before "minute" to avoid "分钟s"
	timeUnitsOrdered := []struct {
		eng   string
		trans string
	}{
		{"seconds", T("seconds")},
		{"minutes", T("minutes")},
		{"hours", T("hours")},
		{"months", T("months")},
		{"weeks", T("weeks")},
		{"years", T("years")},
		{"days", T("days")},
		{"second", T("second")},
		{"minute", T("minute")},
		{"month", T("month")},
		{"hour", T("hour")},
		{"week", T("week")},
		{"year", T("year")},
		{"day", T("day")},
	}
	
	result := duration
	
	// First replace full words (e.g., "5 seconds", "29 minutes")
	// Must be done in order from longest to shortest
	for _, unit := range timeUnitsOrdered {
		if contains(result, " "+unit.eng) {
			result = replace(result, " "+unit.eng, " "+unit.trans)
		}
	}
	
	// Then replace all short form units (handles compound formats like "42m38s")
	for i := 0; i < len(result); i++ {
		for short, trans := range shortForms {
			// Check if current position has the short form
			if i < len(result) && result[i:i+len(short)] == short {
				// Check if the character before is a digit (or it's at the start)
				if i > 0 {
					prevChar := result[i-1]
					if prevChar >= '0' && prevChar <= '9' {
						// Replace the short form
						result = result[:i] + trans + result[i+len(short):]
						i += len(trans) - 1 // Adjust index for the inserted translation
						break
					}
				}
			}
		}
	}
	
	return result
}

// TranslateErrorMessage translates error messages by replacing known error patterns
func TranslateErrorMessage(errMsg string) string {
	if currentLanguage == "en" {
		return errMsg
	}
	
	// Get the translations map for current language
	var translations map[string]string
	switch currentLanguage {
	case "zh_cn":
		translations = ZhCnTranslations
	case "zh_tw":
		translations = ZhTwTranslations
	case "ja":
		translations = JaTranslations
	case "ko":
		translations = KoTranslations
	case "ru":
		translations = RuTranslations
	default:
		return errMsg
	}
	
	// Try to translate the entire message first
	if translated, exists := translations[errMsg]; exists {
		return translated
	}
	
	// If not found, try to replace known error patterns
	result := errMsg
	for key, translation := range translations {
		// Only replace if the key is found in the error message
		if len(key) > 0 && contains(result, key) {
			result = replace(result, key, translation)
		}
	}
	
	return result
}

// contains checks if a string contains a substring (case-sensitive)
func contains(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// indexOf returns the index of the first instance of substr in s, or -1 if substr is not present
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// replace replaces all occurrences of old with new in s
func replace(s, old, new string) string {
	if old == "" {
		return s
	}
	result := ""
	for {
		i := indexOf(s, old)
		if i == -1 {
			result += s
			break
		}
		result += s[:i] + new
		s = s[i+len(old):]
	}
	return result
}
