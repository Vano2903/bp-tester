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
			//log.Println(event)
			switch event.Type {
			case "container":
				// c.l.Debug("container event:", event)
				// c.l.Debug("container checking:", containerID)
				//log.Println("container event\n\n")
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
					// log.Println("[EVENT] Container died:", event.Actor.ID)
					// case "health_status":
					// 	log.Println("[EVENT] Container health status:", event.Actor.ID)
					// case "kill":
					// 	log.Println("[EVENT] Container killed:", event.Actor.ID)
					// case "update":
					// 	log.Println("[EVENT] Container updated:", event.Actor.ID)
				}
			}
		case err := <-errChan:
			c.l.Errorf("error in event handler: %s", err)
		}
	}
}
