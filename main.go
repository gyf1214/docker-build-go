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

func main() {
	flag.Parse()

	info, err := core.GetPackageInfo(*path, *cmd, *deps)
	if err != nil {
		panic(err)
	}

	builder, err := core.NewImageBuild(info)
	if err != nil {
		panic(err)
	}

	err = builder.Build()
	if err != nil {
		panic(err)
	}
}
