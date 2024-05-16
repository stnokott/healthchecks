//go:build integration

// Package healthchecks is a wrapper around the healthchecks.io endpoints.
package healthchecks

import (
	"context"
	"net/http"
	"testing"
)

func TestProject(t *testing.T) {
	config := configFromEnv()

	type args struct {
		ctx  context.Context
		slug string
		msg  string
	}
	tests := []struct {
		name    string
		p       *Project
		args    args
		wantErr bool
	}{
		{
			name: "ping key valid, slug valid",
			p: &Project{
				pingKey: config.PingKey,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				ctx:  context.Background(),
				slug: config.Slug,
			},
			wantErr: false,
		},
		{
			name: "ping key valid, slug invalid",
			p: &Project{
				pingKey: config.PingKey,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				ctx:  context.Background(),
				slug: _slugInvalid,
			},
			wantErr: true,
		},
		{
			name: "ping key invalid",
			p: &Project{
				pingKey: _pingKeyInvalid,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				ctx:  context.Background(),
				slug: config.Slug,
			},
			wantErr: true,
		},
		{
			name: "with body",
			p: &Project{
				pingKey: config.PingKey,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				ctx:  context.Background(),
				slug: config.Slug,
				msg:  "Fuzz Buzz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Start(tt.args.ctx, tt.args.slug); (err != nil) != tt.wantErr {
				t.Errorf("Project.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.p.Success(tt.args.ctx, tt.args.slug); (err != nil) != tt.wantErr {
				t.Errorf("Project.Success() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.p.Fail(tt.args.ctx, tt.args.slug); (err != nil) != tt.wantErr {
				t.Errorf("Project.Fail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.msg != "" {
				if err := tt.p.Log(tt.args.ctx, tt.args.slug, tt.args.msg); (err != nil) != tt.wantErr {
					t.Errorf("Project.Log(%s) error = %v, wantErr %v", tt.args.msg, err, tt.wantErr)
				}
			}
		})
	}
}

func TestProjectSlug(t *testing.T) {
	config := configFromEnv()

	type args struct {
		ctx  context.Context
		slug string
		msg  string
	}
	tests := []struct {
		name    string
		p       *Project
		args    args
		wantErr bool
	}{
		{
			name: "ping key valid, slug valid",
			p: &Project{
				pingKey: config.PingKey,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				ctx:  context.Background(),
				slug: config.Slug,
			},
			wantErr: false,
		},
		{
			name: "ping key valid, slug invalid",
			p: &Project{
				pingKey: config.PingKey,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				ctx:  context.Background(),
				slug: _slugInvalid,
			},
			wantErr: true,
		},
		{
			name: "ping key invalid",
			p: &Project{
				pingKey: _pingKeyInvalid,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				ctx:  context.Background(),
				slug: config.Slug,
			},
			wantErr: true,
		},
		{
			name: "with body",
			p: &Project{
				pingKey: config.PingKey,
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
			},
			args: args{
				ctx:  context.Background(),
				slug: config.Slug,
				msg:  "Fuzz Buzz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notif := tt.p.Slug(tt.args.slug)
			if err := notif.Start(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Project.Slug(%s).Start() error = %v, wantErr %v", tt.args.slug, err, tt.wantErr)
			}
			if err := notif.Success(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Project.Slug(%s).Success() error = %v, wantErr %v", tt.args.slug, err, tt.wantErr)
			}
			if err := notif.Fail(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Project.Slug(%s).Fail() error = %v, wantErr %v", tt.args.slug, err, tt.wantErr)
			}
			if tt.args.msg != "" {
				if err := notif.Log(tt.args.ctx, tt.args.msg); (err != nil) != tt.wantErr {
					t.Errorf("Project.Slug(%s, %s).Log() error = %v, wantErr %v", tt.args.slug, tt.args.msg, err, tt.wantErr)
				}
			}
		})
	}
}

func TestFromURL(t *testing.T) {
	config := configFromEnv()

	type args struct {
		url  string
		opts []Option
	}
	tests := []struct {
		name                string
		useServerPathPrefix bool
		args                args
		wantErrCreate       bool
		wantErrRequest      bool
	}{
		{
			name: "invalid URL",
			args: args{
				url: "\n",
			},
			wantErrCreate: true,
		},
		{
			name: "valid UUID",
			args: args{
				url: config.URLPrefix + "/" + config.UUID,
			},
			wantErrCreate:  false,
			wantErrRequest: false,
		},
		{
			name: "not found",
			args: args{
				url: config.URLPrefix + "/foo",
			},
			wantErrCreate:  false,
			wantErrRequest: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier, err := FromURL(tt.args.url, tt.args.opts...)
			if (err != nil) != tt.wantErrCreate {
				t.Errorf("NewURL(%s) error = %v, wantErr %v", tt.args.url, err, tt.wantErrCreate)
				return
			}
			if err != nil {
				return
			}

			if err := notifier.Success(context.Background()); (err != nil) != tt.wantErrRequest {
				t.Errorf("NewURL(%s).Success() error = %v, wantErr %v", tt.args.url, err, tt.wantErrRequest)
				return
			}
		})
	}
}
