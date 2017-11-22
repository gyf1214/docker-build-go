package core

import "github.com/docker/docker/client"

var docker *client.Client

func init() {
	var err error
	docker, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}
