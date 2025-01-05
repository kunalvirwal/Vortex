package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/image"
)

func ListImages() ([]string, error) {
	images, err := cli.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		return nil, err
	}
	var imageNames []string
	for _, img := range images {
		if len(img.RepoTags) > 0 {
			imageNames = append(imageNames, img.RepoTags[0])
		}
	}
	return imageNames, nil
}

func FindOrPullImage(imageName string) error {
	images, err := ListImages()
	if err != nil {
		return err
	}
	image := strings.Split(imageName, "/")

	imageWithOutRepoName := image[len(image)-1]

	for _, img := range images {
		if img == imageName || img == imageWithOutRepoName {
			return nil
		}
	}
	return PullImage(imageName)
}

func PullImage(imageName string) error {
	fmt.Println("Pulling image", imageName)
	out, err := cli.ImagePull(context.Background(), imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	return nil
}
