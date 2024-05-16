// Package healthchecks is a wrapper around the healthchecks.io endpoints.
package healthchecks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProject(t *testing.T) {
	type args struct {
		ctx  context.Context
		slug string
		msg  string
	}
	tests := []struct {
		name             string
		p                *Project
		serverPathPrefix string
		args             args
		wantErr          bool
	}{
		{
			name: "ping key valid, slug valid",
			p: &Project{
				pingKey: _pingValid,
				opts:    &options{HTTPClient: http.DefaultClient},
			},
			serverPathPrefix: "",
			args: args{
				ctx:  context.Background(),
				slug: _slugValid,
			},
			wantErr: false,
		},
		{
			name: "ping key valid, slug invalid",
			p: &Project{
				pingKey: _pingValid,
				opts:    &options{HTTPClient: http.DefaultClient},
			},
			serverPathPrefix: "",
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
				opts:    &options{HTTPClient: http.DefaultClient},
			},
			serverPathPrefix: "",
			args: args{
				ctx:  context.Background(),
				slug: _slugValid,
			},
			wantErr: true,
		},
		{
			name: "valid with path prefix",
			p: &Project{
				pingKey: _pingValid,
				opts:    &options{HTTPClient: http.DefaultClient},
			},
			serverPathPrefix: "/prefix",
			args: args{
				ctx:  context.Background(),
				slug: _slugValid,
			},
			wantErr: false,
		},
		{
			name: "with body",
			p: &Project{
				pingKey: _pingValid,
				opts:    &options{HTTPClient: http.DefaultClient},
			},
			serverPathPrefix: "",
			args: args{
				ctx:  context.Background(),
				slug: _slugValid,
				msg:  "Fuzz Buzz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(newMockMux(tt.serverPathPrefix))
			defer server.Close()

			tt.p.opts.RootURL = mustURL(server.URL + tt.serverPathPrefix)

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
	type args struct {
		ctx  context.Context
		slug string
		msg  string
	}
	tests := []struct {
		name             string
		p                *Project
		serverPathPrefix string
		args             args
		wantErr          bool
	}{
		{
			name: "ping key valid, slug valid",
			p: &Project{
				pingKey: _pingValid,
				opts:    &options{HTTPClient: http.DefaultClient},
			},
			serverPathPrefix: "",
			args: args{
				ctx:  context.Background(),
				slug: _slugValid,
			},
			wantErr: false,
		},
		{
			name: "ping key valid, slug invalid",
			p: &Project{
				pingKey: _pingValid,
				opts:    &options{HTTPClient: http.DefaultClient},
			},
			serverPathPrefix: "",
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
				opts:    &options{HTTPClient: http.DefaultClient},
			},
			serverPathPrefix: "",
			args: args{
				ctx:  context.Background(),
				slug: _slugValid,
			},
			wantErr: true,
		},
		{
			name: "valid with path prefix",
			p: &Project{
				pingKey: _pingValid,
				opts:    &options{HTTPClient: http.DefaultClient},
			},
			serverPathPrefix: "/prefix",
			args: args{
				ctx:  context.Background(),
				slug: _slugValid,
			},
			wantErr: false,
		},
		{
			name: "witn body",
			p: &Project{
				pingKey: _pingValid,
				opts:    &options{HTTPClient: http.DefaultClient},
			},
			serverPathPrefix: "",
			args: args{
				ctx:  context.Background(),
				slug: _slugValid,
				msg:  "Fuzz Buzz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(newMockMux(tt.serverPathPrefix))
			defer server.Close()

			tt.p.opts.RootURL = mustURL(server.URL + tt.serverPathPrefix)

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
	server := httptest.NewServer(newMockMux(""))
	defer server.Close()
	serverPathPrefix := "/prefix"
	serverWithPrefix := httptest.NewServer(newMockMux(serverPathPrefix))
	defer serverWithPrefix.Close()

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
				url: server.URL + "/" + _uuidValid,
			},
			wantErrCreate:  false,
			wantErrRequest: false,
		},
		{
			name: "valid with path prefix",
			args: args{
				url: serverWithPrefix.URL + serverPathPrefix + "/" + _uuidValid,
			},
			wantErrCreate:  false,
			wantErrRequest: false,
		},
		{
			name: "not found",
			args: args{
				url: server.URL + "/foo",
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
