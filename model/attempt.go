package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Attempt struct {
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CreatedAt             time.Time     `json:"createdAt"`
	ID                    uint          `gorm:"primarykey" json:"-"`
	Code                  string        `json:"code"`
	Output                string        `json:"output,omitempty"`
	Status                AttemptStatus `json:"status"`
	FileContent           []byte        `json:"-"` //for now we remove it as the attempts are fully public
	Executions            []*Execution  `json:"executions" gorm:"-"`
	Best                  *Execution    `json:"best,omitempty" gorm:"-"`
	AverageDuration       time.Duration `json:"averageDuration,omitempty" gorm:"-"`
	AverageDurationString string        `json:"averageDurationString,omitempty" gorm:"-"`
}

func (a Attempt) String() string {
	if a.Executions == nil {
		return fmt.Sprintf("Attempt{ID=%d Code=%q Status=%s}", a.ID, a.Code, a.Status)
	} else {
		return fmt.Sprintf("Attempt{ID=%d Code=%q Status=%s Executions=%s Best=%s Average=%v}", a.ID, a.Code, a.Status, a.Executions, a.Best, a.AverageDuration)
	}
}

type AttemptStatus string

func (a AttemptStatus) String() string {
	return string(a)
}

const (
	AttemptStatusPending AttemptStatus = "pending"
	// AttemptStatusInvalid     AttemptStatus = "invalid"
	AttemptStatusBuilding    AttemptStatus = "building"
	AttemptStatusBuildFailed AttemptStatus = "build_failed"
	AttemptStatusStopped     AttemptStatus = "stopped"
	AttemptStatusSuccess     AttemptStatus = "success"
	AttemptStatusFailed      AttemptStatus = "failed"
)
