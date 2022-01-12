package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
  out, err := cli.ImagePull(ctx, "nginx", types.ImagePullOptions{})
	if err != nil {
    panic(err)
	}
	defer out.Close()

	// block until the image is fully downloaded
	// TODO: Theres probably a better way to do this
	io.Copy(os.Stdout, out)
}
