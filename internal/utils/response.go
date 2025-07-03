package utils

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, statusCode int, errorMessage string) error {
	return c.JSON(statusCode, Response{
		Success: false,
		Error:   errorMessage,
	})
}

func PaginatedResponse(c echo.Context, statusCode int, data interface{}, total int64, page int, limit int) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Data: struct {
			Items      interface{} `json:"items"`
			Total      int64       `json:"total"`
			Page       int         `json:"page"`
			Limit      int         `json:"limit"`
			TotalPages int         `json:"total_pages"`
		}{
			Items:      data,
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: calculateTotalPages(total, limit),
		},
	})
}

func calculateTotalPages(total int64, limit int) int {
	if limit == 0 {
		return 0
	}
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	return totalPages
}