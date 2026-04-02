package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NoCache 是一个 Gin 中间件，用于禁止客户端缓存 HTTP 请求的返回结果.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Cors 是一个 Gin 中间件，用于处理 CORS 请求.
func Cors(c *gin.Context) {
	// 处理预检请求
	if c.Request.Method == http.MethodOptions {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
		return
	}
	c.Next() // 继续处理请求
}

// Secure 是一个 Gin 中间件，用于添加安全相关的 HTTP 头.
func Secure(c *gin.Context) {
	// 控制跨域访问权限
	c.Header("Access-Control-Allow-Origin", "*")
	// 防止页面被嵌入到其他网站的<iframe>中
	// c.Header("X-Frame-Options", "DENY")
	// 禁用浏览器对响应内容的 MIME 类型嗅探
	c.Header("X-Content-Type-Options", "nosniff")
	// 启用浏览器的 XSS 过滤器，检测并拦截跨站脚本攻击
	c.Header("X-XSS-Protection", "1; mode=block")
	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
	c.Next()
}
