package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vano2903/bp-tester/controller"
)

func (h *httpHandler) RefreshTokens(c echo.Context) error {
	ctx := c.Request().Context()
	//get refresh token from cookie
	refreshCookie, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return respError(c, 401, "refresh token cookie not found")
		}
	}
	if refreshCookie.Valid() != nil {
		return respError(c, 401, "refresh token cookie is not valid")
	}
	//generate token pair from token
	accessToken, refreshToken, err := h.controller.GenerateTokenPairFromRefreshToken(ctx, refreshCookie.Value)
	if err != nil {
		h.l.Errorf("error generating a new token pair from a refresh token: %v", err)
		switch err {
		case controller.ErrTokenExpired:
			return respError(c, 401, "refresh token expired", "the refresh token privided is expired you need to login again to generate a new valid refresh token")
		}
		return respError(c, 500, "internal server error")
	}

	//delete refresh token
	if err := h.controller.DeleteRefreshToken(ctx, refreshCookie.Value); err != nil {
		h.l.Errorf("error removing refresh token: %v", err)
		return respError(c, 500, "internal server errro")
	}

	//send tokens as http only cookies
	accessTokenCookie := &http.Cookie{
		Name:     accessTokenCookieName,
		Value:    accessToken.Token,
		Expires:  accessToken.ExpiresAt,
		HttpOnly: true,
	}

	refreshTokenCookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refreshToken.Token,
		Expires:  refreshToken.ExpiresAt,
		HttpOnly: true,
	}

	c.SetCookie(accessTokenCookie)
	c.SetCookie(refreshTokenCookie)

	return respSuccess(c, 200, "token refreshed correctly")
}
