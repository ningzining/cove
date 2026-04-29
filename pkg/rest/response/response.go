package response

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/ningzining/cove/pkg/xerr"
	"github.com/rs/zerolog/log"
)

func PageOk(c *gin.Context, data interface{}, total int64) {
	OK(c, pageData{
		List:  data,
		Total: total,
	})
}

func OK(c *gin.Context, data interface{}) {
	r := Default.Clone()
	r.SetCode(http.StatusOK)
	r.SetMsg("success")
	r.SetData(data)
	c.Set("response", r)
	c.JSON(http.StatusOK, r)
}

func Error(c *gin.Context, err error) {
	code, templateData := xerr.Decode(err)
	var msg string
	if m := t(c, strconv.Itoa(code), templateData); m != "" {
		msg = m
	}
	r := Default.Clone()
	r.SetCode(code)
	r.SetMsg(msg)
	c.Set("response", r)
	c.JSON(http.StatusOK, r)
}

func t(c *gin.Context, messageID string, templateData map[string]interface{}) string {
	localizer, exists := c.Get("localizer")
	if !exists {
		return "" // 容错处理
	}

	msg, err := localizer.(*i18n.Localizer).Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	if err != nil {
		log.Error().Err(err).Msgf("Localize messageID %s, err %s", messageID, err.Error())
		// 如果找不到翻译，返回 ID 或者处理错误
		return ""
	}
	return msg
}
