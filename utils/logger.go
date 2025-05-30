package utils

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Logger is the global logger instance
var Logger *logrus.Logger

// LogFields represents structured log fields
type LogFields map[string]interface{}

// RequestContext holds request-specific information
type RequestContext struct {
	RequestID string
	UserID    string
	Method    string
	Path      string
	IP        string
}

// InitLogger initializes the global logger
func InitLogger() {
	Logger = logrus.New()

	// Set output to stdout (can be redirected to cloud logging)
	Logger.SetOutput(os.Stdout)

	// Set log format to JSON for better parsing in cloud services
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Set log level based on environment
	env := os.Getenv("ENV")
	if env == "production" {
		Logger.SetLevel(logrus.InfoLevel)
	} else {
		Logger.SetLevel(logrus.DebugLevel)
	}

	Logger.Info("Logger initialized successfully")
}

// LogInfo logs info level messages with optional fields
func LogInfo(message string, fields LogFields) {
	if fields != nil {
		Logger.WithFields(logrus.Fields(fields)).Info(message)
	} else {
		Logger.Info(message)
	}
}

// LogError logs error level messages with optional fields
func LogError(message string, err error, fields LogFields) {
	logFields := logrus.Fields{}
	if fields != nil {
		logFields = logrus.Fields(fields)
	}
	if err != nil {
		logFields["error"] = err.Error()
	}
	Logger.WithFields(logFields).Error(message)
}

// LogWarn logs warning level messages with optional fields
func LogWarn(message string, fields LogFields) {
	if fields != nil {
		Logger.WithFields(logrus.Fields(fields)).Warn(message)
	} else {
		Logger.Warn(message)
	}
}

// LogDebug logs debug level messages with optional fields
func LogDebug(message string, fields LogFields) {
	if fields != nil {
		Logger.WithFields(logrus.Fields(fields)).Debug(message)
	} else {
		Logger.Debug(message)
	}
}

// LogWithContext logs with request context information
func LogWithContext(ctx context.Context, level logrus.Level, message string, fields LogFields) {
	logFields := logrus.Fields{}
	if fields != nil {
		logFields = logrus.Fields(fields)
	}

	// Add context information if available
	if reqCtx, ok := ctx.Value("request_context").(*RequestContext); ok {
		logFields["request_id"] = reqCtx.RequestID
		logFields["user_id"] = reqCtx.UserID
		logFields["method"] = reqCtx.Method
		logFields["path"] = reqCtx.Path
		logFields["ip"] = reqCtx.IP
	}

	Logger.WithFields(logFields).Log(level, message)
}

// GetRequestContext extracts request context from Gin context
func GetRequestContext(c *gin.Context) *RequestContext {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
		c.Header("X-Request-ID", requestID)
	}

	userID := c.GetHeader("X-User-ID")

	return &RequestContext{
		RequestID: requestID,
		UserID:    userID,
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		IP:        c.ClientIP(),
	}
}

// LogRequest logs HTTP request information
func LogRequest(c *gin.Context, message string, fields LogFields) {
	reqCtx := GetRequestContext(c)
	ctx := context.WithValue(context.Background(), "request_context", reqCtx)
	LogWithContext(ctx, logrus.InfoLevel, message, fields)
}

// LogRequestError logs HTTP request errors
func LogRequestError(c *gin.Context, message string, err error, fields LogFields) {
	reqCtx := GetRequestContext(c)
	ctx := context.WithValue(context.Background(), "request_context", reqCtx)

	if fields == nil {
		fields = LogFields{}
	}
	if err != nil {
		fields["error"] = err.Error()
	}

	LogWithContext(ctx, logrus.ErrorLevel, message, fields)
}
