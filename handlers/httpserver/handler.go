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
	h.e.Static("/static", "static")

	pages := h.e.Group("")
	pages.GET("/", h.IndexPage)
	pages.GET("/attempt", h.NewAttemptPage)
	pages.GET("/attempt/:code", h.AttemptInfoPage)

	api := h.e.Group("/api/v1")
	attempt := api.Group("/attempt")
	attempt.POST("/new", h.Upload)
	attempt.GET("/info/:code", h.GetAttemptInfo)

	leaderboard := api.Group("/leaderboard")
	leaderboard.GET("/list", h.GetLeaderboard)
}

func (h *httpHandler) IndexPage(c echo.Context) error {
	return c.File("pages/index.html")
}

func (h *httpHandler) NewAttemptPage(c echo.Context) error {
	return c.File("pages/attempt.html")
}

func (h *httpHandler) AttemptInfoPage(c echo.Context) error {
	return c.File("pages/attemptInfo.html")
}
