package main

import (
	"context"
	"log"

	health "github.com/stnokott/healthchecks"
)

const (
	_pingKey = "mysecretpingkey"
	_slug    = "foo"
)

func main() {
	// create a new project using a ping key
	project, err := health.NewProject(_pingKey)
	if err != nil {
		log.Fatalln(err)
	}

	// directly ping using a slug
	err = project.Success(context.TODO(), _slug)
	// ...

	// create notifier instance for a slug
	notifier := project.Slug(_slug)
	err = notifier.Success(context.TODO())
	// ...
}
