package controller

import "errors"

var (
	ErrEmtpySource   = errors.New("empty source")
	ErrSourceTooLong = errors.New("source too long")

	ErrQeueuFull = errors.New("queue is full")

	ErrImageNotCompiledCorrectly = errors.New("image not compiled correctly")

	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")

	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
