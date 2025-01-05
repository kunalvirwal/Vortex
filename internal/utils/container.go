package utils

import (
	"fmt"
	"strings"

	"github.com/kunalvirwal/Vortex/internal/state"
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
			fmt.Println(newParts)
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
	if isNewSalt(salt) {
		return salt
	}
	return generateRandomString(length)

}

func isNewSalt(salt string) bool {
	for _, container := range state.VortexContainers {
		l := len(container.Name)
		if container.Name[l-5:] == salt {
			return false
		}
	}
	return true
}
