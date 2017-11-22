package core

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// BuildRunner run the container & copy back the file
type BuildRunner struct {
	image     string
	build     string
	output    string
	container string
}

// NewBuildRunner return a BuildRunner based on ImageBuilder
func NewBuildRunner(b *ImageBuilder) *BuildRunner {
	return &BuildRunner{
		image:  b.dockerImage,
		build:  b.dockerBuild,
		output: b.pkg.Output,
	}
}

// Clean stop and remove the container used
func (b *BuildRunner) Clean() {
	docker.ContainerRemove(context.Background(), b.container, types.ContainerRemoveOptions{
		Force: true,
	})
}

// Run run the container to build and copy back
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

	tr := tar.NewReader(reader)
	_, err = tr.Next()
	if err != nil {
		return err
	}

	file, err := os.Create(b.output)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, tr)
	if err != nil {
		return err
	}

	return nil
}
