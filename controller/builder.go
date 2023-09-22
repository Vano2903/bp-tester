package controller

import (
	"context"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/vano2903/bp-tester/model"
)

// todo: needs to be checked by the routine monitor
// todo: at start it needs to read all pending attempts in the db and add them to the queue
func (c *Controller) ProcessBuildQueue(ctx context.Context, errChan chan error, ID int, routineMonitor chan int) {
	c.l.Info("starting build queue")
	defer func() {
		if r := recover(); r != nil {
			c.l.Errorf("build queue panic, recovering: \nerror: %v\n\nstack: %s", r, string(debug.Stack()))
		}
		if ctx.Err() == nil {
			routineMonitor <- ID
		} else {
			c.l.Info("build queue not restarting, context was canceled")
		}
	}()

	var currentAttempt *model.Attempt
	for {
		select {
		case <-ctx.Done():
			c.l.Info("context done, stopping build queue")
			if currentAttempt != nil {
				currentAttempt.Status = model.AttemptStatusPending
				if err := c.attemptRepo.UpdateOne(ctx, currentAttempt); err != nil {
					c.l.Errorf("error updating attempt %s: %s", currentAttempt.Code, err)
					errChan <- err
				}
			}
			//gracefully stop the current attempt and wait for it to be cleared
			//the current execution will be lost but the attempt will be in pending state
			//so it will run again when the server starts again
			for {
				time.Sleep(100 * time.Millisecond)
				if currentAttempt == nil {
					break
				}
			}
			return
		case attempt := <-c.buildQueue:
			c.l.Infof("processing %s", attempt.Code)
			currentAttempt = attempt
			if err := c.BuildAttempt(ctx, attempt); err != nil {
				c.l.Errorf("error building attempt %s: %s", attempt.Code, err)
				errChan <- err
			}
			currentAttempt = nil
		}
	}
}

func (c *Controller) BuildAttempt(ctx context.Context, attempt *model.Attempt) error {
	/*
		1. x create image that builds the binary
		2. x check output of build
		2.1 x if build failed - set status to build failed
		3. x create testing container from image
		4. x execute test
		5. check output of test
		5.1 if test failed - set status to failed
		5.2 x if test took more than 2 minutes - set status to timeout
		5.3 check exit status code (if 0 then passed, else failed)
		6. set status to success
		7. get execution time
		8. append execution to attempt
		9. destroy container
		10. repeat from 3 until all tests are done
		11. remove image
		12. clean build dir

		if a test fails keep running the next test but set the status to failed
		if a test times out keep running the next test but set the status to failed
		if a test times out the execution time is 0
	*/
	attempt.Status = model.AttemptStatusBuilding
	if err := c.attemptRepo.UpdateOne(ctx, attempt); err != nil {
		c.l.Errorf("error updating attempt %s: %s", attempt.Code, err)
		return err
	}

	imageName := c.CreateImageName(attempt.Code)

	defer c.CleanBuildDir(ctx)

	c.l.Infof("creating build context for %s", imageName)
	buildContext, err := c.CreateBuildContext(ctx, attempt.FileContent)
	if err != nil {
		return err
	}
	defer buildContext.Close()

	c.l.Infof("creating image %s", imageName)
	image, err := c.CreateImage(ctx, buildContext, imageName, attempt.Code)
	if err != nil {
		return err
	}

	c.l.Infof("checking build status")
	if !c.CheckIfImageCompiled(image.BuildOutput) {
		c.l.Info("build failed")
		attempt.Status = model.AttemptStatusBuildFailed
		//TODO: format build output to be more readable
		attempt.Output = string(image.BuildOutput)
		if err := c.attemptRepo.UpdateOne(ctx, attempt); err != nil {
			c.l.Errorf("error updating attempt %s: %s", attempt.Code, err)
			return err
		}
		return nil
	}

	for i := 1; i <= c.config.Test.RepeatFor; i++ {
		c.l.Infof("running execution number: %d/%d for attempt %s", i, c.config.Test.RepeatFor, attempt.Code)
		//create container
		containerName := c.CreateContainerName(attempt.Code)
		c.l.Infof("creating container %s", containerName)
		containerID, err := c.CreateNewContainer(ctx, containerName, imageName)
		if err != nil {
			return err
		}

		c.l.Infof("starting container %s", containerName)
		execution, err := c.NewExecution(ctx, attempt.ID, i)
		if err != nil {
			c.l.Errorf("error creating execution for attempt %s: %s", attempt.Code, err)
			return err
		}

		containerCtx, cancel := context.WithTimeout(ctx, c.config.Test.Timeout)
		defer cancel()

		//listen for container events
		var isExecutionDone bool
		events := make(chan *model.Event)
		done := make(chan struct{})
		go c.ListenForEvents(containerCtx, containerID, events)

		//read container events
		go func() {
			for {
				select {
				case event := <-events:
					switch event.Action {
					case "die":
						c.l.Infof("container %s died", containerName)
						isExecutionDone = true
						exitCode, err := strconv.Atoi(event.ExitCode)
						if err != nil {
							c.l.Errorf("error converting exit code %s: %s", event.ExitCode, err)
							return
						}
						execution.ExitCode = exitCode
						if exitCode == 0 {
							execution.Status = model.ExecutionStatusPassed
						} else {
							execution.Status = model.ExecutionStatusRunError
						}
						done <- struct{}{}
						return
					}
				case <-containerCtx.Done():
					return
				}
			}
		}()
		//stop the container if the context is done
		go func() {
			<-containerCtx.Done()
			if ctx.Err() != nil {
				c.l.Infof("removing container %s, main context done", containerName)
				if err := c.RemoveContainerByID(ctx, containerID); err != nil {
					c.l.Errorf("error removing container %s: %s", containerName, err)
				}
				return
			}
			if isExecutionDone {
				return
			}
			c.l.Infof("%s timeout reached, removing container %s", c.config.Test.Timeout, containerName)
			execution.Status = model.ExecutionStatusTimeout
			execution.Duration = c.config.Test.Timeout
			execution.ExitCode = -1
			attempt.Executions = append(attempt.Executions, execution)
			attempt.Status = model.AttemptStatusFailed
			if err := c.RemoveContainerByID(ctx, containerID); err != nil {
				c.l.Errorf("error removing container %s: %s", containerName, err)
			}
		}()

		if err := c.StartContainerByID(ctx, containerID); err != nil {
			return err
		}
		<-done

		output, err := c.GetExecutionOutput(ctx, containerID)
		if err != nil {
			return err
		}

		output = strings.Trim(output, "\n")
		outputs := strings.Split(output, "\n")
		execution.Output = outputs[0]
		if !c.CheckExecutionOutput(ctx, execution.Output, c.config.Test.ExpectedOutput) {
			execution.Status = model.ExecutionStatusIncorrectOutput
		}

		time := outputs[len(outputs)-1]
		execution.Duration, err = c.GetExecutionTime(ctx, time)
		if err != nil {
			return err
		}

		c.l.Infof("execution %d finished with status %s", i, execution.Status)
		attempt.Executions = append(attempt.Executions, execution)
		//update attempt
		if err := c.executionRepo.UpdateOne(ctx, execution); err != nil {
			c.l.Errorf("error updating execution %d: %s", execution.ID, err)
			return err
		}

		c.l.Infof("removing container %s of attempt number %d", containerName, i)
		if err := c.RemoveContainerByID(ctx, containerID); err != nil {
			c.l.Errorf("error removing container %s: %s", containerName, err)
		}
		cancel()
	}

	attempt.Status = model.AttemptStatusSuccess
	for _, execution := range attempt.Executions {
		if execution.Status != model.ExecutionStatusPassed {
			attempt.Status = model.AttemptStatusFailed
			break
		}
	}

	if err := c.attemptRepo.UpdateOne(ctx, attempt); err != nil {
		c.l.Errorf("error updating attempt %s: %s", attempt.Code, err)
		return err
	}

	c.l.Infof("removing image %s", imageName)
	return c.RemoveImage(ctx, image.ImageID)
}

func (c *Controller) CleanBuildDir(ctx context.Context) error {
	//cleaning build dir
	c.l.Info("cleaning build dir")
	if err := os.RemoveAll(c.config.Test.BuildDir); err != nil {
		c.l.Errorf("error cleaning build dir %s: %s", c.config.Test.BuildDir, err)
		return err
	}
	if err := os.MkdirAll(c.config.Test.BuildDir, 0777); err != nil {
		c.l.Errorf("error recreating build dir %s: %s", c.config.Test.BuildDir, err)
		return err
	}
	return nil
}

// todo: implement
func (c *Controller) GetExecutionTime(ctx context.Context, t string) (time.Duration, error) {
	a, err := strconv.Atoi(t)
	return time.Duration(a) * time.Millisecond, err
}

func (c *Controller) CheckExecutionOutput(ctx context.Context, output string, expectedValue string) bool {
	// c.l.Debugf("checking output: %q = %q", output, expectedValue)
	return output == expectedValue
}
