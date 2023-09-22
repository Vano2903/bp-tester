package controller

import (
	"context"
	"os"

	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/vano2903/bp-tester/config"
	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/repo"
)

type Controller struct {
	buildDir      string
	config        *config.Config
	lastAttemptID int
	attemptRepo   repo.AttemptRepoer
	executionRepo repo.ExecutionRepoer
	buildQueue    chan *model.Attempt
	l             *logrus.Logger
	cli           *client.Client
}

func NewController(l *logrus.Logger, config *config.Config, attempRepo repo.AttemptRepoer, executionRepo repo.ExecutionRepoer, ctx context.Context) (*Controller, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)

	return &Controller{
		l:             l,
		buildDir:      "./build",
		config:        config,
		cli:           cli,
		attemptRepo:   attempRepo,
		executionRepo: executionRepo,
		lastAttemptID: 0, //i know it's 0 by default, but i want to be explicit
		buildQueue:    make(chan *model.Attempt, 100),
	}, nil
}

func (c *Controller) init() error {
	//create build directory
	if err := os.MkdirAll(c.buildDir, 0777); err != nil {
		return err
	}

	return nil
}
