package core

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/jhoonb/archivex"
)

type ImageBuilder struct {
	pkg         PackageInfo
	tmpPath     string
	dockerFile  string
	dockerTar   string
	dockerImage string
	dockerWd    string
	dockerBuild string
}

const (
	tmpFs        = "/tmp"
	tmpPrefix    = "docker-build-go"
	tmpTar       = "docker.tar"
	dockerFile   = "Dockerfile"
	dockerPrefix = "docker-builder-go-"
	dockerGoSrc  = "/go/src"
	template     = `FROM golang

ENV WORKING_DIR %v

RUN mkdir -p ${WORKING_DIR}
COPY . ${WORKING_DIR}
WORKDIR ${WORKING_DIR}

RUN apt-get update &&\
    apt-get install %v -y &&\
    go-wrapper download &&\
    go-wrapper install

CMD ["go", "build", "-o", "%v", "%v"]
`
)

func NewImageBuild(pkg PackageInfo) (*ImageBuilder, error) {
	tmp, err := ioutil.TempDir(tmpFs, tmpPrefix)
	if err != nil {
		return nil, err
	}

	wd := filepath.Join(dockerGoSrc, pkg.Full)

	ret := &ImageBuilder{
		pkg:         pkg,
		tmpPath:     tmp,
		dockerFile:  filepath.Join(tmp, dockerFile),
		dockerTar:   filepath.Join(tmp, tmpTar),
		dockerImage: dockerPrefix + pkg.Short,
		dockerWd:    wd,
		dockerBuild: filepath.Join(wd, pkg.Build),
	}

	err = os.MkdirAll(ret.tmpPath, 0755)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (b *ImageBuilder) generateDockerfile() error {
	file, err := os.Create(b.dockerFile)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, template,
		b.dockerWd, b.pkg.Deps, b.pkg.Build, b.pkg.Cmd)
	if err != nil {
		return err
	}

	return nil
}

func (b *ImageBuilder) generateTar() error {
	tar := new(archivex.TarFile)
	defer tar.Close()

	err := tar.Create(b.dockerTar)
	if err != nil {
		return err
	}

	err = tar.AddAll(b.pkg.Path, false)
	if err != nil {
		return err
	}

	err = tar.AddFile(b.dockerFile)
	if err != nil {
		return err
	}

	return nil
}

func (b *ImageBuilder) build() error {
	tar, err := os.Open(b.dockerTar)
	defer tar.Close()
	if err != nil {
		return err
	}

	response, err := docker.ImageBuild(context.Background(), tar, types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		Tags:           []string{b.dockerImage},
	})
	if err != nil {
		return err
	}

	defer response.Body.Close()
	reader := bufio.NewScanner(response.Body)
	for reader.Scan() {
		line := reader.Text()
		var result map[string]string
		json.Unmarshal([]byte(line), &result)
		if result["error"] != "" {
			return errors.New(result["error"])
		} else if result["stream"] != "" {
			fmt.Print(result["stream"])
		}
	}

	return nil
}

func (b *ImageBuilder) Clean() {
	// delete tmp dir
	os.RemoveAll(b.tmpPath)

	ctx := context.Background()

	// remove build image
	docker.ImageRemove(ctx, b.dockerImage, types.ImageRemoveOptions{})

	// clean unused images
	arg := filters.NewArgs()
	arg.Add("dangling", "1")
	docker.ImagesPrune(ctx, arg)
}

func (b *ImageBuilder) Build() error {
	err := b.generateDockerfile()
	if err != nil {
		return err
	}

	err = b.generateTar()
	if err != nil {
		return err
	}

	err = b.build()
	if err != nil {
		return err
	}

	return nil
}
