package mock

import (
	"context"
	"time"

	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/repo"
)

func NewRefreshTokenRepoer() repo.RefreshTokenRepoer {
	return &RefreshTokenRepoerMock{
		lastID:  0,
		storage: make(map[string]*model.RefreshToken),
	}

}

type RefreshTokenRepoerMock struct {
	lastID  uint
	storage map[string]*model.RefreshToken
}

func (r *RefreshTokenRepoerMock) FindByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	refreshToken := new(model.RefreshToken)
	refreshToken, ok := r.storage[token]
	if !ok {
		return nil, repo.ErrNotFound
	}
	return refreshToken, nil
}

func (r *RefreshTokenRepoerMock) InsertOne(ctx context.Context, token *model.RefreshToken) error {
	r.lastID++
	token.ID = r.lastID
	token.CreatedAt = time.Now()
	r.storage[token.Token] = token
	return nil
}

func (r *RefreshTokenRepoerMock) DeleteOne(ctx context.Context, token *model.RefreshToken) error {
	delete(r.storage, token.Token)
	return nil
}
