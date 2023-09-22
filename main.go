package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/vano2903/bp-tester/config"
	"github.com/vano2903/bp-tester/controller"
	"github.com/vano2903/bp-tester/handlers/httpserver"
	"github.com/vano2903/bp-tester/pkg/logger"
	"github.com/vano2903/bp-tester/repo/sqliteRepo"
)

const (
	gracefulShutdownTimeout = 15 * time.Second

	StartWebServer      = 1
	StartBuildProcessor = 2
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	l := logger.NewLogger("debug", "text")
	l.Debug("initizalized logger")

	l.Debugf("config: %+v", conf)

	defer func(l *logrus.Logger) {
		if r := recover(); r != nil {
			l.Errorf("panic: recover: %v", r)
			l.Errorf("stacktrace from panic: \n%s", string(debug.Stack()))
		}
	}(l)

	if conf.DB.Driver != "sqlite3" {
		l.Fatal("only sqlite3 driver is supported in the current version")
	}

	attemptRepo, err := sqliteRepo.NewAttemptRepo(conf.DB.Path)
	if err != nil {
		l.Fatal("new attempt sqlite:", err)
	}

	executionRepo, err := sqliteRepo.NewExecutionRepo(conf.DB.Path)
	if err != nil {
		l.Fatal("new execution sqlite:", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	c, err := controller.NewController(l, conf, attemptRepo, executionRepo, ctx)
	if err != nil {
		l.Fatal("new controller:", err)
	}

	e := echo.New()
	httpHandler := httpserver.InitRouter(e, l, c, conf)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGTERM)

	var RoutineMonitor = make(chan int, 100)
	RoutineMonitor <- StartWebServer
	RoutineMonitor <- StartBuildProcessor
	var buildErrChan = make(chan error, 100)
	for {
		select {
		case i := <-interrupt:
			l.Info("main - signal: " + i.String())
			l.Info("main - canceling context")
			cancel()
			gracefulTimer := time.Tick(gracefulShutdownTimeout)
			select {
			case <-gracefulTimer:
				l.Info("main - graceful shutdown timeout reached")
				os.Exit(1)
			case <-httpHandler.Done:
				l.Info("main - http server stopped")
				os.Exit(0)
			}
		case err = <-buildErrChan:
			l.Errorf("build error: %v", err)

		// case err = <-rmq.Error:
		// l.Error(fmt.Errorf("rabbitmq: %w", err))
		default:
		}

		select {
		case ID := <-RoutineMonitor:
			l.Infof("Starting Routine: %d", ID)
			switch ID {
			case StartWebServer:
				go httpserver.StartRouter(ctx, httpHandler, conf, ID, RoutineMonitor)
			case StartBuildProcessor:
				go c.ProcessBuildQueue(ctx, buildErrChan, ID, RoutineMonitor)
			default:
			}
		default:
		}

		time.Sleep(10 * time.Millisecond)
	}

}
