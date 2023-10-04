package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/vano2903/bp-tester/controller"
)

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
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
	pages.GET("/register", h.RegisterPage)
	pages.GET("/login", h.LoginPage)
	pages.GET("/attempt", h.NewAttemptPage)
	pages.GET("/attempt/:code", h.AttemptInfoPage)

	api := h.e.Group("/api/v1")
	user := api.Group("/user")
	user.POST("/register", h.Register)
	user.POST("/login", h.Login)

	tokens := api.Group("/tokens")
	tokens.GET("/refresh", h.RefreshTokens)

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
	//if the user is not logged it it can still do an attempt
	//but it will not be considered in the leaderboard
	return c.File("pages/attempt.html")
}

func (h *httpHandler) RegisterPage(c echo.Context) error {
	return c.File("pages/register.html")
}

func (h *httpHandler) LoginPage(c echo.Context) error {
	return c.File("pages/login.html")
}

func (h *httpHandler) AttemptInfoPage(c echo.Context) error {
	return c.File("pages/attemptInfo.html")
}
