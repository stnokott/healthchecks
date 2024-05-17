//go:build integration

package healthchecks

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

const (
	_uuidInvalid    = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	_pingKeyInvalid = "ijklmnop"
	_slugInvalid    = "bar"
)

func TestRequest(t *testing.T) {
	config := configFromEnv()

	type args struct {
		opts *options
		body io.Reader
		path []string
	}
	tests := []struct {
		name             string
		operations       []string
		serverPathPrefix string
		args             args
		wantErr          bool
	}{
		{
			name:             "uuid valid",
			operations:       []string{"", "/start", "/fail"},
			serverPathPrefix: "",
			args: args{
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
				body: nil,
				path: []string{"/", config.UUID},
			},
			wantErr: false,
		},
		{
			name:             "uuid invalid",
			operations:       []string{""},
			serverPathPrefix: "",
			args: args{
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
				body: nil,
				path: []string{"/", _uuidInvalid},
			},
			wantErr: true,
		},
		{
			name:             "ping key valid, slug valid",
			operations:       []string{"", "/start", "/fail"},
			serverPathPrefix: "",
			args: args{
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
				body: nil,
				path: []string{"/", config.PingKey, "/", config.Slug},
			},
			wantErr: false,
		},
		{
			name:             "ping key valid, slug invalid",
			operations:       []string{""},
			serverPathPrefix: "",
			args: args{
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
				body: nil,
				path: []string{"/", config.PingKey, "/", _slugInvalid},
			},
			wantErr: true,
		},
		{
			name:             "ping key invalid",
			operations:       []string{""},
			serverPathPrefix: "",
			args: args{
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
				body: nil,
				path: []string{"/", _pingKeyInvalid, "/", config.Slug},
			},
			wantErr: true,
		},
		{
			name:             "path invalid",
			operations:       []string{""},
			serverPathPrefix: "",
			args: args{
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
				body: nil,
				path: []string{"/", "invalid"},
			},
			wantErr: true,
		},
		{
			name:             "operation invalid",
			operations:       []string{"/bar"},
			serverPathPrefix: "",
			args: args{
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
				body: nil,
				path: []string{"/", config.UUID},
			},
			wantErr: true,
		},
		{
			name:             "with body",
			operations:       []string{"/log"},
			serverPathPrefix: "",
			args: args{
				opts: &options{
					RootURL:    mustURL(config.URLPrefix),
					HTTPClient: http.DefaultClient,
				},
				body: strings.NewReader("Foo Bar"),
				path: []string{"/", config.UUID},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, op := range tt.operations {
				t.Run(op, func(t *testing.T) {
					path := make([]string, len(tt.args.path)+1)
					copy(path, tt.args.path)
					path[len(path)-1] = op

					if err := request(context.Background(), tt.args.opts, tt.args.body, path...); (err != nil) != tt.wantErr {
						t.Errorf("request() error = %v, wantErr %v", err, tt.wantErr)
					}
				})
			}
		})
	}
}
