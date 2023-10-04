package httpserver

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/vano2903/bp-tester/model"
)

func (h *httpHandler) GetUserFromContext(c echo.Context) (*model.User, error) {
	ctx := c.Request().Context()
	//get access token cookie
	accessTokenCookie, err := c.Cookie(accessTokenCookieName)
	if err != nil {
		return nil, err //http.ErrNoCookie
	}
	if accessTokenCookie.Valid() != nil {
		return nil, errors.New("invalid cooki")
	}

	return h.controller.ValidateAccessTokenAndGetUser(ctx, accessTokenCookie.Value)
}
