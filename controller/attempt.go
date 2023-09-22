package controller

import (
	"context"
	"time"

	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/repo"
)

func (c *Controller) newAttemptCode(ctx context.Context) string {
	var code string
	for {
		code = "at-" + RandStringBytes(5)
		c.l.Debug("checking code ", code)
		if _, err := c.attemptRepo.FindByCode(ctx, code); err == repo.ErrNotFound {
			break
		}
	}

	return code
}

func (c *Controller) newAttemptModel(ctx context.Context, source []byte) *model.Attempt {
	attempt := new(model.Attempt)
	attempt.Code = c.newAttemptCode(ctx)
	attempt.CreatedAt = time.Now()
	attempt.Status = model.AttemptStatusPending
	attempt.FileContent = source
	return attempt
}

func (c *Controller) IsValidSource(ctx context.Context, source []byte) error {
	// c.l.Info("source is valid (not implemented)")
	if len(source) == 0 {
		return ErrEmtpySource
	}
	if len(source) > 10485760 {
		return ErrSourceTooLong
	}
	return nil
}

func (c *Controller) InsertAttempt(ctx context.Context, attempt *model.Attempt) error {
	c.l.Infof("inserting attempt %s", attempt.Code)
	return c.attemptRepo.InsertOne(ctx, attempt)
}

func (c *Controller) NewAttempt(ctx context.Context, source []byte) (*model.Attempt, error) {
	if err := c.IsValidSource(ctx, source); err != nil {
		return nil, err
	}

	attempt := c.newAttemptModel(ctx, source)
	c.l.Infof("creating new attempt with code %s and pushing to build queue", attempt.Code)

	if err := c.InsertAttempt(ctx, attempt); err != nil {
		c.l.Infof("error inserting attempt %s: %s", attempt.Code, err)
		return nil, err
	}

	select {
	case c.buildQueue <- attempt:
	default:
		return nil, ErrQeueuFull
	}

	return attempt, nil
}

func (c *Controller) CalculateAttemptStats(ctx context.Context, attempt *model.Attempt) {
	if attempt.Executions == nil || len(attempt.Executions) == 0 {
		return
	}

	var sumDurations time.Duration
	executions := attempt.Executions
	best := executions[0]
	for i := 0; i < len(executions); i++ {
		sumDurations += executions[i].Duration
		if executions[i].Duration < best.Duration {
			best = executions[i]
		}
		executions[i].DurationString = executions[i].Duration.String()
	}
	attempt.Best = best
	attempt.AverageDuration = time.Duration(int(sumDurations) / len(attempt.Executions))
	attempt.AverageDurationString = attempt.AverageDuration.String()
}

func (c *Controller) GetAttemptByCode(ctx context.Context, code string) (*model.Attempt, error) {
	attempt, err := c.attemptRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	executions, err := c.executionRepo.FindByAttemptID(ctx, attempt.ID)
	if err != nil {
		if err != repo.ErrNotFound {
			c.l.Infof("attempt %s has no executions yet", attempt.Code)
			return attempt, err
		}
		return nil, err
	}
	attempt.Executions = executions
	c.CalculateAttemptStats(ctx, attempt)

	return attempt, nil
}
