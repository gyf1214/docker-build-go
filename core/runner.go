package core

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type BuildRunner struct {
	image     string
	build     string
	output    string
	container string
}

func NewBuildRunner(b *ImageBuilder) *BuildRunner {
	return &BuildRunner{
		image:  b.dockerImage,
		build:  b.dockerBuild,
		output: b.pkg.Output,
	}
}

func (b *BuildRunner) Clean() {
	docker.ContainerRemove(context.Background(), b.container, types.ContainerRemoveOptions{
		Force: true,
	})
}

func (b *BuildRunner) Run() error {
	ctx := context.Background()

	fmt.Println("build start...")

	resp, err := docker.ContainerCreate(ctx, &container.Config{
		Image: b.image,
	}, nil, nil, b.image)
	if err != nil {
		return err
	}
	b.container = resp.ID

	docker.ContainerStart(ctx, b.container, types.ContainerStartOptions{})

	statusCh, errCh := docker.ContainerWait(ctx, b.container, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
		// nothing
	}

	fmt.Println("build finished...")

	reader, _, err := docker.CopyFromContainer(ctx, b.container, b.build)
	if err != nil {
		return err
	}
	defer reader.Close()

	file, err := os.Create(b.output)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}
