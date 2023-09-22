package controller

import (
	"context"
	"time"

	"github.com/vano2903/bp-tester/model"
)

func (c *Controller) newExecutionModel(ctx context.Context, attemptID uint, position int) *model.Execution {
	execution := new(model.Execution)
	execution.Status = model.ExecutionStatusRunning
	execution.Position = position
	execution.ExecutedAt = time.Now()
	execution.AttemptID = attemptID
	return execution
}

func (c *Controller) NewExecution(ctx context.Context, attemptID uint, position int) (*model.Execution, error) {
	execution := c.newExecutionModel(ctx, attemptID, position)
	if err := c.executionRepo.InsertOne(ctx, execution); err != nil {
		return nil, err
	}
	return execution, nil
}
