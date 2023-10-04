package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/vano2903/bp-tester/controller"
	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/repo"
)

type (
	AccessPostRequest struct {
		Username string `json:"username" example:"username"`
		Password string `json:"password" example:"password"`
	}

	AccessResponse struct {
		Username     string              `json:"username" example:"username"`
		AccessToken  *model.AccessToken  `json:"access_token"`
		RefreshToken *model.RefreshToken `json:"refresh_token"`
	}
)

func (h *httpHandler) Register(c echo.Context) error {
	req := new(AccessPostRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	user, err := h.controller.CreateUser(ctx, req.Username, req.Password)
	if err != nil {
		h.l.Errorf("failed to create user: %v", err)
		if err == repo.ErrUsernameTaken {
			return respError(c, 400, "username taken")
		}
		return respError(c, 500, "internal server error")
	}
	accessToken, refreshToken, err := h.controller.GenerateTokenPair(ctx, user.ID)
	if err != nil {
		return respError(c, 500, "internal server error")
	}

	resp := &AccessResponse{
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return respSuccess(c, 200, "user created correctly", resp)
}

func (h *httpHandler) Login(c echo.Context) error {
	req := new(AccessPostRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	user, err := h.controller.Login(ctx, req.Username, req.Password)
	if err != nil {
		h.l.Errorf("failed to login user: %v", err)
		switch err {
		case repo.ErrNotFound, controller.ErrInvalidCredentials:
			return respError(c, 400, "invalid credentials")
		}
		return respError(c, 500, "internal server error")
	}
	accessToken, refreshToken, err := h.controller.GenerateTokenPair(ctx, user.ID)
	if err != nil {
		return respError(c, 500, "internal server error")
	}

	resp := &AccessResponse{
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return respSuccess(c, 200, "user logged in correctly", resp)
}
