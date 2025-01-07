package docker

import (
	"context"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/kunalvirwal/Vortex/types"
	"golang.org/x/exp/rand"
)

func GenerateContainerName(cfg *types.ContainerConfig) string {

	imageName := ""
	parts := strings.Split(cfg.Image, "@")
	if len(parts) > 1 {
		imageName = parts[0]
	} else {
		parts = strings.Split(cfg.Image, ":")
		if len(parts) > 1 {
			newParts := strings.Split(parts[0], "/")
			if len(newParts) > 1 {
				imageName = newParts[1]
			} else {
				imageName = parts[0]
			}
		} else {
			newParts := strings.Split(cfg.Image, "/")
			if len(newParts) > 1 {
				imageName = newParts[1]
			} else {
				imageName = parts[0]
			}
		}
	}
	name := imageName + "-" + cfg.Service + "-" + generateRandomString(5)
	return name
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	salt := string(b)
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return "unable to generate random string"
	}
	if isNewSalt(salt, length, containers) {
		return salt
	}
	return generateRandomString(length)

}

func isNewSalt(salt string, length int, containers []dockertypes.Container) bool {
	for _, container := range containers {
		for _, name := range container.Names {
			l := len(name)
			if l >= length && name[l-length:] == salt {
				return false
			}
		}
	}
	return true
}
