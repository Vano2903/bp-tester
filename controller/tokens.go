package controller

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/repo"
)

// todo: get time from config
func (c *Controller) CreateRefreshToken(ctx context.Context, userID uint) (*model.RefreshToken, error) {
	refreshTokenDuration := time.Hour * 24 * 7

	ran, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	refreshTokenValue := ran.String()
	refreshToken := new(model.RefreshToken)
	refreshToken.Token = refreshTokenValue
	refreshToken.ExpiresAt = time.Now().Add(refreshTokenDuration)
	refreshToken.UserID = userID

	return refreshToken, nil
}

// todo: get duration from config
func (c *Controller) GenerateAccessToken(ctx context.Context, userID uint) (*model.AccessToken, error) {
	accessTokenDuration := time.Hour

	ran, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	accessTokenValue := ran.String()
	accessToken := new(model.AccessToken)
	accessToken.Token = accessTokenValue
	accessToken.ExpiresAt = time.Now().Add(accessTokenDuration)
	accessToken.UserID = userID

	return accessToken, nil
}

func (c *Controller) GenerateTokenPair(ctx context.Context, userID uint) (*model.AccessToken, *model.RefreshToken, error) {
	accessToken, err := c.GenerateAccessToken(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := c.CreateRefreshToken(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	err = c.accessTokenRepo.InsertOne(ctx, accessToken)
	if err != nil {
		return nil, nil, err
	}

	err = c.refreshTokenRepo.InsertOne(ctx, refreshToken)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, err
}

func (c *Controller) IsRefreshTokenExpired(ctx context.Context, refreshToken string) (bool, error) {
	token, err := c.refreshTokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		if err == repo.ErrNotFound {
			return true, nil
		}
		return true, err
	}

	return token.ExpiresAt.Before(time.Now()), nil
}

func (c *Controller) IsAccessTokenExpired(ctx context.Context, accessToken string) (bool, error) {
	token, err := c.accessTokenRepo.FindByToken(ctx, accessToken)
	if err != nil {
		if err == repo.ErrNotFound {
			return true, nil
		}
		return true, err
	}

	return token.ExpiresAt.Before(time.Now()), nil
}

func (c *Controller) GenerateTokenPairFromRefreshToken(ctx context.Context, refreshToken string) (*model.AccessToken, *model.RefreshToken, error) {
	refreshTokenExpired, err := c.IsRefreshTokenExpired(ctx, refreshToken)
	if err != nil {
		return nil, nil, err
	}
	if refreshTokenExpired {
		return nil, nil, ErrTokenExpired
	}

	token := new(model.RefreshToken)
	token.Token = refreshToken
	if err = c.refreshTokenRepo.DeleteOne(ctx, token); err != nil {
		return nil, nil, err
	}
	return c.GenerateTokenPair(ctx, token.UserID)
}

func (c *Controller) ValidateAccessTokenAndGetUser(ctx context.Context, accessToken string) (*model.User, error) {
	token, err := c.accessTokenRepo.FindByToken(ctx, accessToken)
	if err != nil {
		if err == repo.ErrNotFound {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	accessTokenExpired, err := c.IsAccessTokenExpired(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if accessTokenExpired {
		return nil, ErrTokenExpired
	}

	return c.userRepo.FindByID(ctx, token.UserID)
}

//TODO: Implement revoke token
// the refreshToken can easly be deleted but for the access token
// we either store it in a blacklist or we store all the access tokens in the database.
// if we store all the tokens then we can also implement a "logout from all devices" feature
// while if we do it with the blacklist we need to not allow
// all the tokens issued before the "logout from all devices" to be valid (which is not a bad thing)

func (c *Controller) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	return c.refreshTokenRepo.DeleteOne(ctx, &model.RefreshToken{
		Token: refreshToken,
	})
}
