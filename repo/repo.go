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
	InsertOne(ctx context.Context, attempt *model.Attempt) error
	UpdateOne(ctx context.Context, attempt *model.Attempt) error
}

type ExecutionRepoer interface {
	FindByID(ctx context.Context, id uint) (*model.Execution, error)
	FindByAttemptID(ctx context.Context, attemptID uint) ([]*model.Execution, error)
	InsertOne(ctx context.Context, execution *model.Execution) error
	UpdateOne(ctx context.Context, execution *model.Execution) error
}

var (
	ErrNotFound = errors.New("not found")
)
