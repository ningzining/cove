package middleware

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type LoggerConfig struct {
	Out io.Writer
}

func Logger() gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerToFile 日志记录到文件
func LoggerWithConfig(conf LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		var body string
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete:
			bf := bytes.NewBuffer(nil)
			wt := bufio.NewWriter(bf)
			_, err := io.Copy(wt, c.Request.Body)
			if err != nil {
				log.Error().Err(err).Msg("copy body error")
				err = nil
			}
			rb, _ := io.ReadAll(bf)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(rb))
			body = string(rb)
			// 过滤敏感数据
			body = filterSensitiveData(body)
			if len(body) > 1<<20 {
				body = body[:1<<20] + "......[TRUNCATED]"
			}
		}

		c.Next()

		if c.Request.Method == http.MethodOptions {
			return
		}
		// 结束时间
		endTime := time.Now()
		rt, bl := c.Get("response")
		var result = ""
		if bl {
			rb, err := json.Marshal(rt)
			if err != nil {
				log.Error().Err(err).Msg("json Marshal result error")
			} else {
				result = string(rb)
				result = filterSensitiveData(result)
				if len(result) > 1<<20 {
					result = result[:1<<20] + "......[TRUNCATED]"
				}
			}
		}
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		if raw != "" {
			path = path + "?" + raw
		}
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		param := LogFormatterParams{
			Request:    c.Request,
			Body:       body,
			Response:   result,
			TimeStamp:  endTime,
			StatusCode: statusCode,
			Latency:    latencyTime,
			ClientIP:   clientIP,
			Method:     reqMethod,
			Path:       path,
		}
		log.Info().Msg(LogFormatter(param))
	}
}

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	Request  *http.Request
	Body     string
	Response string
	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var LogFormatter = func(param LogFormatterParams) string {
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[Logger] %v | %3d | %13v | %15s | %-7s %#v\n%s",
		param.TimeStamp.Format(time.DateTime),
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}

// filterSensitiveData 过滤JSON字符串中的敏感数据
func filterSensitiveData(jsonStr string) string {
	if jsonStr == "" {
		return ""
	}

	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		// 如果不是有效的JSON，直接返回原字符串
		return jsonStr
	}

	// 递归处理数据
	data = filterData(data)

	// 重新编码为JSON字符串
	result, err := json.Marshal(data)
	if err != nil {
		return jsonStr
	}

	return string(result)
}

// filterData 递归过滤数据中的敏感信息
func filterData(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		// 处理JSON对象
		for key, value := range v {
			// 不区分大小写检查字段名是否包含password关键词
			if strings.Contains(strings.ToLower(key), "password") {
				// 替换密码值
				v[key] = "******"
			} else {
				// 递归处理子数据
				v[key] = filterData(value)
			}
		}
	case []interface{}:
		// 处理JSON数组
		for i, item := range v {
			v[i] = filterData(item)
		}
	}
	return data
}
