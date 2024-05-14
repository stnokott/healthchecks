package healthchecks

import (
	"reflect"
	"testing"
)

func TestUrlOption(t *testing.T) {
	type args struct {
		opts *options
	}
	tests := []struct {
		name     string
		opt      urlOption
		args     args
		wantOpts *options
		wantErr  bool
	}{
		{
			name:     "empty",
			opt:      "",
			args:     args{opts: &options{RootURL: nil}},
			wantOpts: &options{RootURL: nil},
			wantErr:  true,
		},
		{
			name:     "invalid",
			opt:      "foo",
			args:     args{opts: &options{RootURL: nil}},
			wantOpts: &options{RootURL: nil},
			wantErr:  true,
		},
		{
			name:     "valid",
			opt:      "https://example.com",
			args:     args{opts: &options{RootURL: nil}},
			wantOpts: &options{RootURL: mustURL("https://example.com")},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.opt.apply(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("urlOption.apply() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.opts, tt.wantOpts) {
				t.Errorf("urlOption.apply() result mismatch:\ngot =  %#v\nwant = %#v", tt.args.opts, tt.wantOpts)
			}
		})
	}
}
