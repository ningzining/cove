package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIdKey = "X-Request-Id"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		requestId := c.GetHeader(requestIdKey)
		if requestId == "" {
			requestId = c.GetHeader(strings.ToLower(requestIdKey))
		}
		if requestId == "" {
			requestId = uuid.New().String()
		}
		// 将 RequestID 保存到 HTTP Header头中，供后续链路使用，Header 的键为 `X-Request-Id`
		c.Request.Header.Set(requestIdKey, requestId)
		// 将 RequestID 保存到 context.Context 中，以便后续程序使用
		c.Set(requestIdKey, requestId)
		// 继续处理请求
		c.Next()
	}
}
