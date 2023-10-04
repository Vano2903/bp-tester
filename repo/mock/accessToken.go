package mock

import (
	"context"
	"time"

	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/repo"
)

func NewAccessTokenRepoer() repo.AccessTokenRepoer {
	return &AccessTokenRepoerMock{
		lastID:  0,
		storage: make(map[string]*model.AccessToken),
	}

}

type AccessTokenRepoerMock struct {
	lastID  uint
	storage map[string]*model.AccessToken
}

func (r *AccessTokenRepoerMock) FindByToken(ctx context.Context, token string) (*model.AccessToken, error) {
	refreshToken := new(model.AccessToken)
	refreshToken, ok := r.storage[token]
	if !ok {
		return nil, repo.ErrNotFound
	}
	return refreshToken, nil
}

func (r *AccessTokenRepoerMock) InsertOne(ctx context.Context, token *model.AccessToken) error {
	r.lastID++
	token.ID = r.lastID
	token.CreatedAt = time.Now()
	r.storage[token.Token] = token
	return nil
}

func (r *AccessTokenRepoerMock) DeleteOne(ctx context.Context, token *model.AccessToken) error {
	delete(r.storage, token.Token)
	return nil
}
