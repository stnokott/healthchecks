package healthchecks

import (
	"fmt"
	"net/url"
	"time"
)

// TODO: HTTP client
type options struct {
	RootURL *url.URL
	Timeout time.Duration
}

func defaultOptions() *options {
	return &options{
		RootURL: mustURL("https://hc-ping.com"),
		Timeout: 10 * time.Second,
	}
}

func optsFromDefaults(opts []Option) (*options, error) {
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

var _ Option = urlOption("")

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
// The format should be http[s]://example.com[/suffix].
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

var _ Option = timeoutOption(0)

type timeoutOption time.Duration

func (t timeoutOption) apply(opts *options) error {
	if t < 0 {
		return fmt.Errorf("timeout is %d, needs to be > 0", t)
	}
	opts.Timeout = time.Duration(t)
	return nil
}

// WithTimeout sets the timeout for signalling requests (/start, ...).
//
// The default is 10s.
func WithTimeout(t time.Duration) Option {
	return timeoutOption(t)
}
