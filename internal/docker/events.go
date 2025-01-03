package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/events"
)

// listen to docker events
func ListenEvents() {

	ctx := context.Background()
	eventChan, errChan := cli.Events(ctx, events.ListOptions{})

	for {
		select {
		case event := <-eventChan:
			HandleEvent(event)
		case err := <-errChan:
			fmt.Println("Error in listening to docker events:", err)
		}
	}
}

func HandleEvent(event events.Message) {
	fmt.Println("Event:", event)
	if event.Type == events.ContainerEventType {
		switch event.Action {
		case "die":
			// go tracker.ContainerDied(event.Actor.ID)
		case "start":
			// go tracker.ContainerStarted(event.Actor.ID)
		case "stop":
			// go tracker.ContainerStopped(event.Actor.ID)
		}
	}
}
