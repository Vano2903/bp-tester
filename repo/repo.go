package repo

import (
	"context"
	"errors"

	"github.com/vano2903/bp-tester/model"
)

// interface
type AttemptRepoer interface {
	FindByID(ctx context.Context, id uint) (*model.Attempt, error)
	FindByCode(ctx context.Context, code string) (*model.Attempt, error)
	FindByStatus(ctx context.Context, statuses ...model.AttemptStatus) ([]*model.Attempt, error)
	InsertOne(ctx context.Context, attempt *model.Attempt) error
	UpdateOne(ctx context.Context, attempt *model.Attempt) error
}

type ExecutionRepoer interface {
	FindByID(ctx context.Context, id uint) (*model.Execution, error)
	FindByAttemptID(ctx context.Context, attemptID uint) ([]*model.Execution, error)
	InsertOne(ctx context.Context, execution *model.Execution) error
	UpdateOne(ctx context.Context, execution *model.Execution) error
}

type UserRepoer interface {
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	InsertOne(ctx context.Context, user *model.User) error
}

type AccessTokenRepoer interface {
	FindByToken(ctx context.Context, token string) (*model.AccessToken, error)
	InsertOne(ctx context.Context, token *model.AccessToken) error
	DeleteOne(ctx context.Context, token *model.AccessToken) error
}

type RefreshTokenRepoer interface {
	FindByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	InsertOne(ctx context.Context, token *model.RefreshToken) error
	DeleteOne(ctx context.Context, token *model.RefreshToken) error
}

var (
	ErrNotFound      = errors.New("not found")
	ErrUsernameTaken = errors.New("username taken")
)
