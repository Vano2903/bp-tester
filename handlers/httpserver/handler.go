package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/vano2903/bp-tester/controller"
)

type (
	httpHandler struct {
		e          *echo.Echo
		controller *controller.Controller
		l          *logrus.Logger
		Done       chan struct{}
	}
)

func NewHttpHandler(e *echo.Echo, c *controller.Controller, l *logrus.Logger) *httpHandler {
	return &httpHandler{
		e:          e,
		controller: c,
		l:          l,
		Done:       make(chan struct{}),
	}
}

// Registers only the routes and links functions
func (h *httpHandler) RegisterRoutes() {
	api := h.e.Group("/api/v1")
	attempt := api.Group("/attempt")
	attempt.POST("/new", h.Upload)
	attempt.GET("/info/:code", h.GetAttemptInfo)
}
