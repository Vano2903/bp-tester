package mock

import (
	"context"

	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/repo"
)

func NewExecutionRepoer() repo.ExecutionRepoer {
	return &ExecutionMockRepo{
		storage: make(map[uint]*model.Execution),
	}
}

type ExecutionMockRepo struct {
	lastID  uint
	storage map[uint]*model.Execution
}

func (r *ExecutionMockRepo) FindByID(ctx context.Context, id uint) (*model.Execution, error) {
	execution, ok := r.storage[id]
	if !ok {
		return nil, repo.ErrNotFound
	}
	return execution, nil
}

func (r *ExecutionMockRepo) FindByAttemptID(ctx context.Context, attemptID uint) ([]*model.Execution, error) {
	var executions []*model.Execution
	for _, e := range r.storage {
		if e.AttemptID == attemptID {
			executions = append(executions, e)
		}
	}
	return executions, nil
}

func (r *ExecutionMockRepo) InsertOne(ctx context.Context, attempt *model.Execution) error {
	r.lastID++
	attempt.ID = r.lastID
	r.storage[attempt.ID] = attempt
	return nil
}

func (r *ExecutionMockRepo) UpdateOne(ctx context.Context, attempt *model.Execution) error {
	r.storage[attempt.ID] = attempt
	return nil
}
