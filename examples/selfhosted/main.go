package main

import (
	"context"
	"log"

	health "github.com/stnokott/healthchecks"
)

const (
	_uuid          = "12345678-abcd-1234-5678-999999999999"
	_selfHostedURL = "https://example.com"
)

func main() {
	// create a new check using a UUID
	check, err := health.NewUUID(_uuid, health.WithURL(_selfHostedURL))
	if err != nil {
		log.Fatalln(err)
	}

	// ping
	err = check.Success(context.TODO())
	// ...
}
