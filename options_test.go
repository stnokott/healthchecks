package healthchecks

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestWithURL(t *testing.T) {
	type args struct {
		opts *options
	}
	tests := []struct {
		name     string
		url      string
		args     args
		wantOpts *options
		wantErr  bool
	}{
		{
			name:     "empty",
			url:      "",
			args:     args{opts: &options{RootURL: nil}},
			wantOpts: &options{RootURL: nil},
			wantErr:  true,
		},
		{
			name:     "invalid",
			url:      "\n",
			args:     args{opts: &options{RootURL: nil}},
			wantOpts: &options{RootURL: nil},
			wantErr:  true,
		},
		{
			name:     "valid",
			url:      "https://example.com",
			args:     args{opts: &options{RootURL: nil}},
			wantOpts: &options{RootURL: mustURL("https://example.com")},
			wantErr:  false,
		},
		{
			name:     "with path",
			url:      "https://example.com/health",
			args:     args{opts: &options{RootURL: nil}},
			wantOpts: &options{RootURL: mustURL("https://example.com/health")},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithURL(tt.url)

			if err := option.apply(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("WithURL(%s).apply() error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.opts, tt.wantOpts) {
				t.Errorf("WithURL(%s).apply() result mismatch:\ngot =  %#v\nwant = %#v", tt.url, tt.args.opts, tt.wantOpts)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	type args struct {
		opts *options
	}
	tests := []struct {
		name     string
		timeout  time.Duration
		args     args
		wantOpts *options
		wantErr  bool
	}{
		{
			name:    "valid",
			timeout: 5 * time.Second,
			args: args{opts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       10 * time.Second,
			}}},
			wantOpts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       5 * time.Second,
			}},
			wantErr: false,
		},
		{
			name:    "zero",
			timeout: 0 * time.Second,
			args: args{opts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       10 * time.Second,
			}}},
			wantOpts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       0 * time.Second,
			}},
			wantErr: false,
		},
		{
			name:    "negative",
			timeout: -1 * time.Second,
			args: args{opts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       10 * time.Second,
			}}},
			wantOpts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       10 * time.Second,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithTimeout(tt.timeout)

			if err := option.apply(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("WithTimeout(%s).apply() error = %v, wantErr %v", tt.timeout, err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.opts, tt.wantOpts) {
				t.Errorf("WithTimeout(%s).apply() result mismatch:\ngot =  %#v\nwant = %#v", tt.timeout, tt.args.opts, tt.wantOpts)
			}
		})
	}
}

func TestWithHTTPClient(t *testing.T) {
	type args struct {
		opts *options
	}
	tests := []struct {
		name     string
		client   *http.Client
		args     args
		wantOpts *options
		wantErr  bool
	}{
		{
			name: "valid",
			client: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       999 * time.Second,
			},
			args: args{opts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       10 * time.Second,
			}}},
			wantOpts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       999 * time.Second,
			}},
			wantErr: false,
		},
		{
			name: "partially nil",
			client: &http.Client{
				Transport: http.DefaultTransport,
				Timeout:   123 * time.Second,
			},
			args: args{opts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       10 * time.Second,
			}}},
			wantOpts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: nil,
				Jar:           nil,
				Timeout:       123 * time.Second,
			}},
			wantErr: false,
		},
		{
			name:   "nil",
			client: nil,
			args: args{opts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       10 * time.Second,
			}}},
			wantOpts: &options{HTTPClient: &http.Client{
				Transport:     http.DefaultTransport,
				CheckRedirect: http.DefaultClient.CheckRedirect,
				Jar:           http.DefaultClient.Jar,
				Timeout:       10 * time.Second,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithHTTPClient(tt.client)

			if err := option.apply(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("WithHTTPClient().apply() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.opts, tt.wantOpts) {
				t.Errorf("WithHTTPClient().apply() result mismatch:\ngot =  %#v\nwant = %#v", tt.args.opts, tt.wantOpts)
			}
		})
	}
}
