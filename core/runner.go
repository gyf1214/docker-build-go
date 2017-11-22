package core

type BuildRunner struct {
	image     string
	output    string
	container string
}

func NewBuildRunner(b *ImageBuilder) *BuildRunner {
	return &BuildRunner{
		image:  b.dockerImage,
		output: b.pkg.Output,
	}
}

func (b *BuildRunner) Run() error {
	// resp, err := docker.ContainerCreate(context.Background(), &container.Config{
	// 	Image: b.image,
	// }, nil, nil, b.image)
	return nil
}
