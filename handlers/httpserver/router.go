package httpserver

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/vano2903/bp-tester/config"
	"github.com/vano2903/bp-tester/controller"
)

func InitRouter(e *echo.Echo, l *logrus.Logger, controller *controller.Controller, conf *config.Config) *httpHandler {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRemoteIP: true,
		LogMethod:   true,
		LogURI:      true,
		LogStatus:   true,
		LogLatency:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			l.Infof("ip=%q method=%q uri=%q status=%d latency=%q",
				v.RemoteIP,
				v.Method,
				v.URI,
				v.Status,
				v.Latency)
			return nil
		},
	}))
	// e.Use(middleware.Recover())
	// e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
	// 	StackSize:         4 << 10, // 4 KB
	// 	DisableStackAll:   true,
	// 	DisablePrintStack: false,
	// 	LogLevel:          1,
	// }))

	// e.Use(middleware.SecureWithConfig(middleware.DefaultSecureConfig))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	echo.NotFoundHandler = func(c echo.Context) error {
		return respError(c, http.StatusNotFound, "invalid endpoint")
	}

	e.HideBanner = true
	e.HidePort = true

	httpHandler := NewHttpHandler(e, controller, l)
	httpHandler.RegisterRoutes()
	return httpHandler
}

func StartRouter(ctx context.Context, h *httpHandler, conf *config.Config, ID int, RoutineMonitor chan int) {
	defer func() {
		if r := recover(); r != nil {
			h.l.Errorf("router panic, recovering: \nerror: %v\n\nstack: %s", r, string(debug.Stack()))
		}
		if ctx.Err() == nil {
			h.l.Info("router recovered, restarting")
			RoutineMonitor <- ID
		} else {
			h.l.Info("API not restarting, context was canceled")
		}
	}()

	go func() {
		<-ctx.Done()
		if err := h.e.Shutdown(ctx); err != nil {
			h.l.Errorf("error stopping http server: %v", err)
		}
	}()

	addr := "0.0.0.0:" + conf.HTTP.Port
	h.l.Infof(">>> STARTING API ON: %s >>>", addr)
	// var err error
	err := h.e.Start(addr)
	if err != nil {
		if ctx.Err() != nil {
			h.l.Info(">>> STOPPING SERVER >>>")
			h.Done <- struct{}{}
		} else if err == http.ErrServerClosed {
			h.l.Info(">>> SERVER CLOSED >>>")
		} else {
			h.l.Errorf(">>> ERROR STARTING SERVER: %v", err.Error())
		}
	}
}
