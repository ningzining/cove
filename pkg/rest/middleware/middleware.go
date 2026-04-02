package middleware

import "github.com/gin-gonic/gin"

func Setup(r *gin.Engine) {
	r.Use(Recovery()).
		Use(RequestId()).
		Use(Logger()).
		Use(NoCache).
		Use(Cors).
		Use(Secure).
		Use(I18n())
}
