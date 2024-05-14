package main

import (
	"context"
	"log"

	health "github.com/stnokott/healthchecks"
)

const _uuid = "12345678-abcd-1234-5678-999999999999"

func main() {
	// create a new check using a UUID
	check, err := health.NewUUID(_uuid)
	if err != nil {
		log.Fatalln(err)
	}

	// ping
	err = check.Success(context.TODO())
	// ...
}
