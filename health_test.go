package healthchecks

import (
	"reflect"
	"testing"
)

func TestNewProject(t *testing.T) {
	type args struct {
		pingKey string
		opts    []Option
	}
	tests := []struct {
		name    string
		args    args
		want    *Project
		wantErr bool
	}{
		{
			name: "no options",
			args: args{
				pingKey: "ping key",
				opts:    []Option{},
			},
			want: &Project{
				pingKey: "ping key",
				opts:    defaultOptions(),
			},
			wantErr: false,
		},
		{
			name: "empty ping key",
			args: args{
				pingKey: "",
				opts:    []Option{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no ping key",
			args: args{
				pingKey: "",
				opts:    []Option{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with option",
			args: args{
				pingKey: "foo bar",
				opts: []Option{
					WithURL("https://example.com"),
				},
			},
			want: &Project{
				pingKey: "foo bar",
				opts: &options{
					RootURL:    mustURL("https://example.com"),
					HTTPClient: defaultOptions().HTTPClient,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProject(tt.args.pingKey, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProject() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestProjectSlug(t *testing.T) {
	type args struct {
		slug string
	}
	tests := []struct {
		name string
		p    *Project
		args args
		want Notifier
	}{
		{
			name: "valid",
			p: &Project{
				pingKey: "fooBar",
				opts:    defaultOptions(),
			},
			args: args{slug: "sluggySlug"},
			want: &Check{
				path: "fooBar/sluggySlug",
				opts: defaultOptions(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Slug(tt.args.slug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Project.Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUUID(t *testing.T) {
	type args struct {
		uuid string
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    *Check
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				uuid: "abc-def",
				opts: []Option{},
			},
			want: &Check{
				path: "/abc-def",
				opts: defaultOptions(),
			},
			wantErr: false,
		},
		{
			name: "empty uuid",
			args: args{
				uuid: "",
				opts: []Option{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUUID(tt.args.uuid, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromURL(t *testing.T) {
	type args struct {
		u    string
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    Notifier
		wantErr bool
	}{
		{
			name: "simple URL",
			args: args{
				u:    "https://example.com/foo-bar-123",
				opts: []Option{},
			},
			want: &Check{
				path: "/foo-bar-123",
				opts: &options{
					RootURL:    mustURL("https://example.com"),
					HTTPClient: defaultOptions().HTTPClient,
				},
			},
			wantErr: false,
		},
		{
			name: "subpaths",
			args: args{
				u:    "https://example.com/fuzz/foo-bar-123",
				opts: []Option{},
			},
			want: &Check{
				path: "/fuzz/foo-bar-123",
				opts: &options{
					RootURL:    mustURL("https://example.com"),
					HTTPClient: defaultOptions().HTTPClient,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromURL(tt.args.u, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
