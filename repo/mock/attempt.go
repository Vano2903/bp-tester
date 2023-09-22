package mock

import (
	"context"
	"time"

	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/repo"
)

func NewAttemptRepoer() repo.AttemptRepoer {
	return &AttemptRepoerMock{
		storage: make(map[string]*model.Attempt),
	}
}

type AttemptRepoerMock struct {
	lastID  uint
	storage map[string]*model.Attempt
}

func (r *AttemptRepoerMock) FindByID(ctx context.Context, id uint) (*model.Attempt, error) {
	for _, entity := range r.storage {
		if entity.ID == id {
			return entity, nil
		}
	}
	return nil, repo.ErrNotFound
}

func (r *AttemptRepoerMock) FindByCode(ctx context.Context, code string) (*model.Attempt, error) {
	for _, entity := range r.storage {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, repo.ErrNotFound
}

func (r *AttemptRepoerMock) InsertOne(ctx context.Context, attempt *model.Attempt) error {
	r.lastID++
	attempt.ID = r.lastID
	attempt.CreatedAt = time.Now()
	r.storage[attempt.Code] = attempt
	return nil
}

func (r *AttemptRepoerMock) UpdateOne(ctx context.Context, attempt *model.Attempt) error {
	r.storage[attempt.Code] = attempt
	return nil
}
