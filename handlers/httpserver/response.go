package httpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type HttpError struct {
	Code    int    `json:"code" example:"400"`
	IsError bool   `json:"is_error" example:"true"`
	Message string `json:"message" example:"Bad Request"`
	Details string `json:"details,omitempty" example:"Bad Request With More Info"`
}

func respError(c echo.Context, code int, message, details string) error {
	h := HttpError{
		IsError: true,
		Code:    code,
		Message: message,
		Details: details,
	}

	return c.JSON(code, h)
}

func respErrorf(c echo.Context, code int, message, details string, args ...string) error {
	h := HttpError{
		IsError: true,
		Code:    code,
		Message: message,
		Details: fmt.Sprintf(details, args),
	}

	return c.JSON(code, h)
}

func respErrorFromHttpError(c echo.Context, err *HttpError) error {
	h := HttpError{
		IsError: true,
		Code:    err.Code,
		Message: err.Message,
		Details: err.Details,
	}
	return c.JSON(err.Code, h)
}

type HttpSuccess struct {
	Code    int         `json:"code" example:"200"`
	IsError bool        `json:"is_error" example:"false"`
	Message string      `json:"message" example:"OK"`
	Data    interface{} `json:"data,omitempty"`
}

func respSuccess(c echo.Context, code int, message string, data ...interface{}) error {
	h := HttpSuccess{
		Code:    code,
		IsError: false,
		Message: message,
	}

	if len(data) > 0 {
		h.Data = data[0]
	}

	return c.JSON(code, h)
}

func respSuccessf(c echo.Context, code int, message string, args ...string) error {
	h := HttpSuccess{
		Code:    code,
		IsError: false,
		Message: fmt.Sprintf(message, args),
	}

	return c.JSON(code, h)
}
