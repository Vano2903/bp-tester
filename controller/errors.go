package controller

import "errors"

var (
	ErrEmtpySource   = errors.New("empty source")
	ErrSourceTooLong = errors.New("source too long")

	ErrQeueuFull     = errors.New("queue is full")

	ErrImageNotCompiledCorrectly = errors.New("image not compiled correctly")
)
