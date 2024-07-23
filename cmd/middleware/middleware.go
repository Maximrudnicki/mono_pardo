package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		startTime := time.Now()

		// Read the request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
		}
		// Restore the request body to the original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		// Process request
		c.Next()

		// Request latency
		latency := time.Since(startTime)

		// Request data
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		clientIP := c.ClientIP()

		// Status code
		statusCode := c.Writer.Status()

		// Read the response body
		responseBody := c.Writer

		// Log format
		logMessage := fmt.Sprintf("[%s] %s %s - %s %d (%dms) %s\nRequest: %s\nResponse: %s\n\n\n",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			query,
			statusCode,
			latency.Milliseconds(),
			clientIP,
			string(requestBody),
			responseBody,
		)

		// Write log to file
		logFile, err := os.OpenFile("logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Error opening log file:", err)
			return
		}
		defer logFile.Close()

		if _, err := logFile.WriteString(logMessage); err != nil {
			log.Println("Error writing log:", err)
			return
		}
	}
}
