package langs

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/isayme/port-detector/locales"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var DefaultLanguage = language.English
var bundle *i18n.Bundle
var localizer *i18n.Localizer

func initBundle() {
	bundle = i18n.NewBundle(DefaultLanguage)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 加载所有翻译文件
	entries, err := locales.LocalsFS.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "active.") {
			filePath := entry.Name()
			data, err := locales.LocalsFS.ReadFile(filePath)
			if err != nil {
				panic(err)
			}

			_, err = bundle.ParseMessageFileBytes(data, filePath)
			if err != nil {
				panic(err)
			}
		}
	}
}

func initLocalizer(lang string) {
	// 如果传入空字符串，使用默认语言
	if lang == "" {
		localizer = i18n.NewLocalizer(bundle, DefaultLanguage.String())
	}

	// 创建包含回退语言的本地化器
	localizer = i18n.NewLocalizer(bundle, lang, DefaultLanguage.String())
}

func Setup() {
	initBundle()

	initLocalizer(DefaultLanguage.String())
}

func Localize(id string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
	})
}
