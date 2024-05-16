//go:build integration

package healthchecks

import (
	"go-simpler.org/env"
)

type testConfig struct {
	UUID      string `env:"TESTING_UUID,required"`
	PingKey   string `env:"TESTING_PING_KEY,required"`
	Slug      string `env:"TESTING_SLUG,required"`
	URLPrefix string `env:"TESTING_URL,required"`
}

func configFromEnv() *testConfig {
	c := new(testConfig)
	if err := env.Load(c, nil); err != nil {
		panic(err)
	}
	return c
}
