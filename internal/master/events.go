package master

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/events"
	"github.com/kunalvirwal/Vortex/internal/docker"
)

// listen to docker events
func ListenEvents() {

	ctx := context.Background()
	cli := docker.NewClient()
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
	// fmt.Println("Event:", event)
	if event.Type == events.ContainerEventType {
		switch event.Action {
		case "die":
			go ContainerDied(event)
		case "start":
			// go tracker.ContainerStarted(event.Actor.ID)

		}
	}
}

// Events in life of a Docker container:
// Event: create
// Event: attach
// Event: connect
// Event: start
// Event: kill
// Event: disconnect
// Event: die
// Event: destroy
