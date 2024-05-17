//go:build integration

package healthchecks

import (
	"context"
	"net/http"
	"testing"
)

var config = configFromEnv()

func TestProjectEndpoints(t *testing.T) {
	type args struct {
		slug string
	}
	tests := []struct {
		name    string
		p       *Project
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			p: &Project{
				pingKey: config.PingKey,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				slug: config.Slug,
			},
			wantErr: false,
		},
		{
			name: "invalid URL",
			p: &Project{
				pingKey: config.PingKey,
				opts: &options{
					RootURL:    mustURL("https://example.com"),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				slug: config.Slug,
			},
			wantErr: true,
		},
		{
			name: "invalid ping key",
			p: &Project{
				pingKey: _pingKeyInvalid,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				slug: config.Slug,
			},
			wantErr: true,
		},
		{
			name: "invalid slug",
			p: &Project{
				pingKey: config.PingKey,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				slug: _slugInvalid,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Start(context.Background(), tt.args.slug); (err != nil) != tt.wantErr {
				t.Errorf("Project.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.p.Fail(context.Background(), tt.args.slug); (err != nil) != tt.wantErr {
				t.Errorf("Project.Fail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.p.Success(context.Background(), tt.args.slug); (err != nil) != tt.wantErr {
				t.Errorf("Project.Success() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectLog(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr bool
	}{
		{
			name:    "valid",
			msg:     "Lorem Ipsum",
			wantErr: false,
		},
		{
			name:    "empty message",
			msg:     "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Project{
				pingKey: config.PingKey,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			}
			if err := p.Log(context.Background(), config.Slug, tt.msg); (err != nil) != tt.wantErr {
				t.Errorf("Project.Log(%s, %s) error = %v, wantErr %v", config.Slug, tt.msg, err, tt.wantErr)
			}
		})
	}
}

func TestProjectExitStatus(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		wantErr bool
	}{
		{
			name:    "negative",
			code:    -1,
			wantErr: true,
		},
		{
			name:    "zero",
			code:    0,
			wantErr: false,
		},
		{
			name:    "positive",
			code:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Project{
				pingKey: config.PingKey,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			}

			if err := p.ExitStatus(context.Background(), config.Slug, tt.code); (err != nil) != tt.wantErr {
				t.Errorf("Project.ExitStatus(%s, %d) error = %v, wantErr %v", config.Slug, tt.code, err, tt.wantErr)
			}
		})
	}
}

func TestCheckEndpoints(t *testing.T) {
	tests := []struct {
		name    string
		c       *Check
		wantErr bool
	}{
		{
			name: "valid",
			c: &Check{
				path: config.UUID,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid URL",
			c: &Check{
				path: config.UUID,
				opts: &options{
					RootURL:    mustURL("https://example.com"),
					HTTPClient: http.DefaultClient,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid uuid",
			c: &Check{
				path: _uuidInvalid,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Start(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Check.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.c.Fail(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Check.Fail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.c.Success(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Check.Success() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckLog(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr bool
	}{
		{
			name:    "valid",
			msg:     "Lorem Ipsum",
			wantErr: false,
		},
		{
			name:    "empty message",
			msg:     "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Check{
				path: config.UUID,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			}
			if err := c.Log(context.Background(), tt.msg); (err != nil) != tt.wantErr {
				t.Errorf("Check.Log(%s) error = %v, wantErr %v", tt.msg, err, tt.wantErr)
			}
		})
	}
}

func TestCheckExitStatus(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		wantErr bool
	}{
		{
			name:    "negative",
			code:    -1,
			wantErr: true,
		},
		{
			name:    "zero",
			code:    0,
			wantErr: false,
		},
		{
			name:    "positive",
			code:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Check{
				path: config.UUID,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			}
			if err := c.ExitStatus(context.Background(), tt.code); (err != nil) != tt.wantErr {
				t.Errorf("Check.ExitStatus(%d) error = %v, wantErr %v", tt.code, err, tt.wantErr)
			}
		})
	}
}
