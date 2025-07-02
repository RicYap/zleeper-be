// internal/utils/logger.go
package utils

import (
	"log"
	"os"
	"time"

	echo "github.com/labstack/echo/v4"
)

type Logger struct {
	*log.Logger
	file *os.File
}

func NewLogger(filename string) *Logger {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	return &Logger{
		Logger: log.New(file, "", log.LstdFlags|log.Lmicroseconds),
		file:   file,
	}
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func RequestLoggerMiddleware(logger *Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Log request details
			logger.Printf("Request started: %s %s", c.Request().Method, c.Request().URL.Path)

			err := next(c)

			// Log response details
			duration := time.Since(start)
			logger.Printf("Request completed: %s %s | Status: %d | Duration: %v",
				c.Request().Method, c.Request().URL.Path, c.Response().Status, duration)

			return err
		}
	}
}