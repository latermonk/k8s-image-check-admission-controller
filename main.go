package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		panic(err)
	}

	// get image information
	imageName := "ibackchina2018/ubuntu-sshd:huge"
	args := filters.NewArgs(filters.Arg("reference", imageName))
	images, err := cli.ImageList(ctx, types.ImageListOptions{Filters: args})
	if err != nil {
		panic(err)
	}

	if len(images) != 1 {
		panic(fmt.Errorf("Can't find image %s", imageName))
	}

	// get image size，shown in MB
	imageSize := float32(images[0].Size) / (1024 * 1024)
	fmt.Printf("IMage %s size =  %.2f MB\n", imageName, imageSize)
}
