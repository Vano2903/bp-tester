package controller

import (
	"context"
	"testing"

	"github.com/vano2903/bp-tester/config"
	"github.com/vano2903/bp-tester/controller"
	"github.com/vano2903/bp-tester/model"
	"github.com/vano2903/bp-tester/pkg/logger"
	"github.com/vano2903/bp-tester/repo/sqliteRepo"
)

var c *controller.Controller

func newController() (*controller.Controller, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	l := logger.NewLogger(cfg.Log.Level, cfg.Log.Type)

	// attemptRepo := mock.NewAttemptRepoer()
	// executionRepo := mock.NewExecutionRepoer()
	attemptRepo, err := sqliteRepo.NewAttemptRepo(cfg.DB.Path)
	if err != nil {
		return nil, err
	}
	executionRepo, err := sqliteRepo.NewExecutionRepo(cfg.DB.Path)
	if err != nil {
		return nil, err
	}
	return controller.NewController(l, cfg, attemptRepo, executionRepo, context.Background())
}

func init() {
	var err error
	c, err = newController()
	if err != nil {
		panic(err)
	}
}

var incorrectOutput = []byte(`package main
import (
	"fmt"
	"time"
)
func main() {
	time.Sleep(5 * time.Second)
	fmt.Println("this output is incorrect")
}`)

var buildErrorSource = []byte(`package main
import (
	"fmt"
	"time"
)
func main() {
	time.Sleep(5 * time.Second // missing closing bracket
	fmt.Println("this code wont compile")
}`)

var validSourceWithCorrectOutput = []byte(`package main
import (
	"fmt"
	"time"
)
func main() {
	time.Sleep(124 * time.Millisecond)
	fmt.Println("800382571")
}`)

type source struct {
	Source                []byte
	Message               string
	ExpectedAttemptStatus model.AttemptStatus
	ExpectedExecution     model.ExecutionStatus
}

func TestLoadNewAttempt(t *testing.T) {
	t.Log("TestLoadNewAttempt")

	attempt, err := c.NewAttempt(context.Background(), incorrectOutput)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(attempt)
}

func TestBuild(t *testing.T) {

	t.Log("TestBuild")
	sources := []source{
		// {incorrectOutput, "incorrect output", model.AttemptStatusFailed, model.ExecutionStatusIncorrectOutput},
		// {buildErrorSource, "build error source", model.AttemptStatusBuildFailed, ""},
		{validSourceWithCorrectOutput, "valid source with correct output", model.AttemptStatusSuccess, model.ExecutionStatusPassed},
	}

	for _, s := range sources {
		t.Log("checking build with ", s.Message)
		// errChan := make(chan error)
		ctx, cancel := context.WithCancel(context.Background())
		// c.ProcessBuildQueue(ctx, errChan)
		defer cancel()

		attempt, err := c.NewAttempt(ctx, s.Source)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("attempt:", attempt)

		err = c.BuildAttempt(ctx, attempt)
		if err != nil {
			t.Fatal(err)
		}

		attempt, err = c.GetAttemptByCode(ctx, attempt.Code)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("attempt:", attempt)
		if attempt.Status != s.ExpectedAttemptStatus {
			t.Fatalf("expected attempt status %s, got %s", s.ExpectedAttemptStatus, attempt.Status)
		}
		if s.ExpectedExecution != "" {
			for _, e := range attempt.Executions {
				if e.Status != s.ExpectedExecution {
					t.Fatalf("expected execution status %s, got %s", s.ExpectedExecution, e.Status)
				}
			}
		}
		cancel()
	}
}
