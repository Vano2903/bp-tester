package httpserver

import (
	"bytes"

	"github.com/labstack/echo/v4"
	"github.com/vano2903/bp-tester/controller"
	"github.com/vano2903/bp-tester/repo"
)

func (h *httpHandler) GetAttemptInfo(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return respError(c, 400, "invalid code", "code can't be empty")
	}
	ctx := c.Request().Context()
	attempt, err := h.controller.GetAttemptByCode(ctx, code)
	if err != nil {
		switch err {
		case repo.ErrNotFound:
			return respErrorf(c, 404, "attempt not found", "there are no attempts with code %s", code)
		}
		h.l.Errorf("error getting attempt by code: %v", err)
		return respError(c, 500, "unexpected error")
	}
	return respSuccess(c, 200, "attempt found", attempt)
}

func (h *httpHandler) Upload(c echo.Context) error {
	ctx := c.Request().Context()
	body := c.Request().Body
	defer body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	h.l.Debugf("body: %s", buf.String())
	attempt, err := h.controller.NewAttempt(ctx, buf.Bytes())
	if err != nil {
		h.l.Errorf("error creating new attempt: %v", err)
		switch err {
		case controller.ErrEmtpySource:
			return respError(c, 400, "emtpy source", "source code can't be emtpy")
		case controller.ErrSourceTooLong:
			return respError(c, 413, "source too long", "source code can't be longer than 10MB")
		case controller.ErrQeueuFull:
			return respError(c, 503, "build queue is full", "the system is processing too many attempts, try again later")
		}
		return respError(c, 500, "unexpected error")
	}

	attempt.Best = nil
	attempt.Executions = nil
	return respSuccess(c, 200, "attempt created", attempt)
}
