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
	buildDir         string
	config           *config.Config
	attemptRepo      repo.AttemptRepoer
	executionRepo    repo.ExecutionRepoer
	userRepo         repo.UserRepoer
	accessTokenRepo  repo.AccessTokenRepoer
	refreshTokenRepo repo.RefreshTokenRepoer
	buildQueue       chan *model.Attempt
	l                *logrus.Logger
	cli              *client.Client
}

func NewController(ctx context.Context,
	l *logrus.Logger,
	config *config.Config,
	attempRepo repo.AttemptRepoer,
	executionRepo repo.ExecutionRepoer,
	userRepo repo.UserRepoer,
	accessTokenRepo repo.AccessTokenRepoer,
	refreshTokenRepo repo.RefreshTokenRepoer) (*Controller, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)

	return &Controller{
		l:                l,
		buildDir:         "./build",
		config:           config,
		cli:              cli,
		attemptRepo:      attempRepo,
		executionRepo:    executionRepo,
		userRepo:         userRepo,
		accessTokenRepo:  accessTokenRepo,
		refreshTokenRepo: refreshTokenRepo,
		buildQueue:       make(chan *model.Attempt, 100),
	}, nil
}

func (c *Controller) init() error {
	if err := os.MkdirAll(c.buildDir, 0777); err != nil {
		return err
	}

	return nil
}
