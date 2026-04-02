package middleware

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

func I18n() gin.HandlerFunc {
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	if _, err := bundle.LoadMessageFile("etc/i18n/zh.json"); err != nil {
		log.Fatal().Err(err).Msg("LoadMessageFile zh.json failed")
	}
	if _, err := bundle.LoadMessageFile("etc/i18n/en.json"); err != nil {
		log.Fatal().Err(err).Msg("LoadMessageFile en.json failed")
	}
	return func(c *gin.Context) {
		// 1. 获取语言标签
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = "zh"
		}
		// 2. 创建 Localizer
		// Localize 接受多个 tag 参数，会按顺序尝试匹配
		localizer := i18n.NewLocalizer(bundle, lang)

		// 3. 将 Localizer 存入 Context，方便后续 Controller 使用
		c.Set("localizer", localizer)

		c.Next()
	}
}
