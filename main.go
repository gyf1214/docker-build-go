package main

import (
	"flag"

	"github.com/gyf1214/docker-build-go/core"
)

var (
	path = flag.String("path", "", "path to build")
	cmd  = flag.String("cmd", "", "command to build relative to package path")
	deps = flag.String("deps", "", "apt-get dependencies")
)

func work() error {
	info, err := core.GetPackageInfo(*path, *cmd, *deps)
	if err != nil {
		return err
	}

	builder, err := core.NewImageBuild(info)
	if err != nil {
		return err
	}

	// defer builder.Clean()
	err = builder.Build()
	if err != nil {
		return err
	}

	runner := core.NewBuildRunner(builder)
	// defer runner.Clean()
	err = runner.Run()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	err := work()
	if err != nil {
		panic(err)
	}
}
