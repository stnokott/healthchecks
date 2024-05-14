// Package healthchecks is a wrapper around the healthchecks.io endpoints.
package healthchecks

import (
	"context"
	"net/http/httptest"
	"testing"
)

func TestProject(t *testing.T) {
	type args struct {
		ctx  context.Context
		slug string
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
				pingKey: _pingValid,
				opts:    &options{},
			},
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
				opts:    &options{},
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
				opts:    &options{},
			},
			args: args{
				ctx:  context.Background(),
				slug: _slugValid,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(newMockMux())
			defer server.Close()

			tt.p.opts.RootURL = mustURL(server.URL)

			if err := tt.p.Start(tt.args.ctx, tt.args.slug); (err != nil) != tt.wantErr {
				t.Errorf("Project.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.p.Success(tt.args.ctx, tt.args.slug); (err != nil) != tt.wantErr {
				t.Errorf("Project.Success() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.p.Fail(tt.args.ctx, tt.args.slug); (err != nil) != tt.wantErr {
				t.Errorf("Project.Fail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectSlug(t *testing.T) {
	type args struct {
		ctx  context.Context
		slug string
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
				pingKey: _pingValid,
				opts:    &options{},
			},
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
				opts:    &options{},
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
				opts:    &options{},
			},
			args: args{
				ctx:  context.Background(),
				slug: _slugValid,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(newMockMux())
			defer server.Close()

			tt.p.opts.RootURL = mustURL(server.URL)

			notif := tt.p.Slug(tt.args.slug)
			if err := notif.Start(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Project.Slug(%s).Start() error = %v, wantErr %v", tt.args.slug, err, tt.wantErr)
			}
			if err := notif.Success(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Project.Slug(%s).Success() error = %v, wantErr %v", tt.args.slug, err, tt.wantErr)
			}
			if err := notif.Fail(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Project.Slug(%s).Fail() error = %v, wantErr %v", err, tt.args.slug, tt.wantErr)
			}
		})
	}
}
