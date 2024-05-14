package healthchecks

import (
	"fmt"
	"net/url"
)

// TODO: timeout
// TODO: HTTP client
type options struct {
	RootURL *url.URL
}

func defaultOptions() *options {
	return &options{
		RootURL: mustURL("https://hc-ping.com"),
	}
}

func fromDefaults(opts []Option) (*options, error) {
	options := defaultOptions()
	for _, o := range opts {
		if err := o.apply(options); err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}
	return options, nil
}

// Option applies a configuration option.
type Option interface {
	apply(opts *options) error
}

type urlOption string

func (u urlOption) apply(opts *options) error {
	parsed, err := url.Parse(string(u))
	if err != nil {
		return fmt.Errorf("parse URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("invalid URL: %s", u)
	}
	opts.RootURL = parsed
	return nil
}

// WithURL overrides the default URL https://hc-ping.com.
//
// Format for u should be (http/https)://(host).
func WithURL(u string) Option {
	return urlOption(u)
}

func mustURL(u string) *url.URL {
	parsed, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return parsed
}
