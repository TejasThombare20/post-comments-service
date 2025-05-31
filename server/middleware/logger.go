package middleware

import (
	"time"

	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/gin-gonic/gin"
)

// Logger returns a gin.HandlerFunc for logging HTTP requests using logrus
func Logger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			// Log using logrus instead of standard output
			fields := utils.LogFields{
				"client_ip":   param.ClientIP,
				"timestamp":   param.TimeStamp.Format(time.RFC3339),
				"method":      param.Method,
				"path":        param.Path,
				"protocol":    param.Request.Proto,
				"status_code": param.StatusCode,
				"latency":     param.Latency.String(),
				"user_agent":  param.Request.UserAgent(),
			}

			if param.ErrorMessage != "" {
				fields["error"] = param.ErrorMessage
				utils.LogError("HTTP Request", nil, fields)
			} else {
				utils.LogInfo("HTTP Request", fields)
			}

			// Return empty string since we're handling logging ourselves
			return ""
		},
		Output: nil, // We handle output through logrus
	})
}
