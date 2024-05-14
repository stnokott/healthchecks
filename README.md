[![Go Reference](https://pkg.go.dev/badge/github.com/stnokott/healthchecks.svg)](https://pkg.go.dev/github.com/stnokott/healthchecks)
![GitHub Release](https://img.shields.io/github/v/release/stnokott/healthchecks)


# `healthchecks` - Wrapper for [healthchecks.io](https://healthchecks.io)

Allows simple signalling of statuses to [healthchecks.io](https://healthchecks.io).

# Example usage

## Using [official](https://healthchecks.io) endpoint

### For projects (a group of multiple checks)

You need:
- your project's **ping key**, generated from your project's settings
- one or multiple check **slugs**

```go
package main

import (
	"context"
	"log"

	health "github.com/stnokott/healthchecks"
)

const (
	_pingKey = "mysecretpingkey"
	_slug = "foo"
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
	notifier := project.Slug(_slug_)
	err = notifier.Success(context.TODO())
	// ...
}
```

### For single checks

You need:
- your check's **UUID**

```go
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
```

## Using [self-hosted](https://healthchecks.io/docs/self_hosted) endpoint.

By default, `https://hc-ping.com` is used as endpoint.

Use the `health.WithURL` option to provide your own host.

```go
// for a project
project, err := health.NewProject(pingKey, health.WithURL("https://example.com"))
// ...

// for a UUID-based check
check, err := health.NewUUID(uuid, health.WithURL("https://example.com"))
// ...
```
