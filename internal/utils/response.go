package utils

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Success 	bool        `json:"success"`
	Message 	string      `json:"message,omitempty"`
	Data    	interface{} `json:"data,omitempty"`
	MetaData 	interface{} `json:"metadata,omitempty"`
	Error   	string      `json:"error,omitempty"`
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

func PaginatedResponse(c echo.Context, statusCode int, data interface{}, metadata interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Data: data,
		MetaData: metadata,
	})
}