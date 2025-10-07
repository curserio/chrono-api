package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
)

type Translator struct {
	translations map[language.Tag]map[string]string
}

func NewTranslator(localesPath string) (*Translator, error) {
	t := &Translator{
		translations: make(map[language.Tag]map[string]string),
	}

	// Load translations for supported languages
	supportedLangs := []language.Tag{language.English, language.Russian}
	for _, lang := range supportedLangs {
		path := filepath.Join(localesPath, fmt.Sprintf("%s.json", lang.String()))
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read translations file %s: %w", path, err)
		}

		var translations map[string]string
		if err := json.Unmarshal(data, &translations); err != nil {
			return nil, fmt.Errorf("failed to parse translations file %s: %w", path, err)
		}

		t.translations[lang] = translations
	}

	return t, nil
}

func (t *Translator) Translate(lang language.Tag, key string) string {
	if translations, ok := t.translations[lang]; ok {
		if translation, ok := translations[key]; ok {
			return translation
		}
	}
	return key // Возвращаем ключ, если перевод не найден
}
