package healthchecks

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

const (
	_uuidValid      = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	_uuidInvalid    = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	_pingValid      = "abcdefgh"
	_pingKeyInvalid = "ijklmnop"
	_slugValid      = "foo"
	_slugInvalid    = "bar"
)

func newMockMux(pathPrefix string) http.Handler {
	r := mux.NewRouter()

	uuidHandler := func(w http.ResponseWriter, r *http.Request) {
		uuid := mux.Vars(r)["uuid"]
		switch uuid {
		case _uuidValid:
			fmt.Fprint(w, "OK")
		case _uuidInvalid:
			// yes, the official healthchecks.io endpoint returns 200 when a UUID doesn't exist.
			fmt.Fprint(w, "OK (not found)")
		default:
			w.WriteHeader(400)
			fmt.Fprint(w, "invalid url format")
		}
	}

	pingKeyHandler := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		pingKey, slug := vars["pingKey"], vars["slug"]
		if pingKey == _pingKeyInvalid || slug == _slugInvalid {
			fmt.Fprint(w, "OK (not found)")
			return
		}
		if pingKey == _pingValid && slug == _slugValid {
			fmt.Fprint(w, "OK")
			return
		}
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid URL format")
	}

	r.HandleFunc(pathPrefix+"/{uuid}", uuidHandler)
	r.HandleFunc(pathPrefix+"/{uuid}/start", uuidHandler)
	r.HandleFunc(pathPrefix+"/{uuid}/fail", uuidHandler)

	r.HandleFunc(pathPrefix+"/{pingKey}/{slug}", pingKeyHandler)
	r.HandleFunc(pathPrefix+"/{pingKey}/{slug}/start", pingKeyHandler)
	r.HandleFunc(pathPrefix+"/{pingKey}/{slug}/fail", pingKeyHandler)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid url format")
	})

	return r
}

func TestRequest(t *testing.T) {
	type args struct {
		opts *options
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
				opts: &options{},
				path: []string{"/", _uuidValid},
			},
			wantErr: false,
		},
		{
			name:             "uuid invalid",
			operations:       []string{""},
			serverPathPrefix: "",
			args: args{
				opts: &options{},
				path: []string{"/", _uuidInvalid},
			},
			wantErr: true,
		},
		{
			name:             "ping key valid, slug valid",
			operations:       []string{"", "/start", "/fail"},
			serverPathPrefix: "",
			args: args{
				opts: &options{},
				path: []string{"/", _pingValid, "/", _slugValid},
			},
			wantErr: false,
		},
		{
			name:             "ping key valid, slug invalid",
			operations:       []string{""},
			serverPathPrefix: "",
			args: args{
				opts: &options{},
				path: []string{"/", _pingValid, "/", _slugInvalid},
			},
			wantErr: true,
		},
		{
			name:             "ping key invalid",
			operations:       []string{""},
			serverPathPrefix: "",
			args: args{
				opts: &options{},
				path: []string{"/", _pingKeyInvalid, "/", _slugValid},
			},
			wantErr: true,
		},
		{
			name:             "path invalid",
			operations:       []string{""},
			serverPathPrefix: "",
			args: args{
				opts: &options{},
				path: []string{"/", "invalid"},
			},
			wantErr: true,
		},
		{
			name:             "operation invalid",
			operations:       []string{"/bar"},
			serverPathPrefix: "",
			args: args{
				opts: &options{},
				path: []string{"/", _uuidValid},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, op := range tt.operations {
				t.Run(op, func(t *testing.T) {
					server := httptest.NewServer(newMockMux(tt.serverPathPrefix))
					defer server.Close()

					tt.args.opts.RootURL = mustURL(server.URL)
					path := make([]string, len(tt.args.path)+1)
					copy(path, tt.args.path)
					path[len(path)-1] = op

					if err := request(context.Background(), tt.args.opts, path...); (err != nil) != tt.wantErr {
						t.Errorf("request() error = %v, wantErr %v", err, tt.wantErr)
					}
				})
			}
		})
	}
}
