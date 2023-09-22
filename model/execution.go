package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Execution struct {
	gorm.Model `json:"-"`
	// ID         uint `gorm:"primarykey" json:"-"`
	// CreatedAt  time.Time
	// UpdatedAt  time.Time       `json:"-"`
	// DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
	AttemptID      uint            `json:"-"`
	Position       int             `json:"position"`
	ExitCode       int             `json:"exitCode"`
	ExecutedAt     time.Time       `json:"executedAt"`
	Duration       time.Duration   `json:"duration"`
	DurationString string          `json:"durationString"`
	Output         string          `json:"output"`
	Status         ExecutionStatus `json:"status"`
}

func (e Execution) String() string {
	return fmt.Sprintf("Execution{ID=%d AttemptID=%d Position=%d ExitCode=%d ExecutedAt=%s Duration=%s Output=%s Status=%s}",
		e.ID, e.AttemptID, e.Position, e.ExitCode, e.ExecutedAt, e.Duration, e.Output, e.Status)
}

type ExecutionStatus string

func (e ExecutionStatus) String() string {
	return string(e)
}

const (
	ExecutionStatusRunning         ExecutionStatus = "running"
	ExecutionStatusRunError        ExecutionStatus = "run_error"
	ExecutionStatusIncorrectOutput ExecutionStatus = "incorrect_output"
	ExecutionStatusTimeout         ExecutionStatus = "timeout"
	ExecutionStatusStopped         ExecutionStatus = "stopped"
	ExecutionStatusFailed          ExecutionStatus = "failed"
	ExecutionStatusPassed          ExecutionStatus = "passed"
)
