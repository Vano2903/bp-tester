package controller

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/vano2903/bp-tester/model"
)

func (c *Controller) ListenForEvents(ctx context.Context, containerID string, eventsChan chan *model.Event) {
	eventChan, errChan := c.cli.Events(ctx, types.EventsOptions{})
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-eventChan:
			switch event.Type {
			case "container":
				// c.l.Debug("container event:", event)
				// c.l.Debug("container checking:", containerID)
				switch event.Action {
				case "die":
					if event.Actor.ID == containerID {
						c.l.Infof("container %s died", containerID)
						eventsChan <- &model.Event{
							Action:      "die",
							ContainerID: event.Actor.ID,
							ExitCode:    event.Actor.Attributes["exitCode"],
						}
					}
				}
			}
		case err := <-errChan:
			c.l.Errorf("error in event handler: %s", err)
		}
	}
}
