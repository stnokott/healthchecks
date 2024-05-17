package healthchecks

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type options struct {
	RootURL    *url.URL
	HTTPClient *http.Client
}

func defaultOptions() *options {
	return &options{
		RootURL: mustURL("https://hc-ping.com"),
		HTTPClient: &http.Client{
			Transport:     http.DefaultTransport,
			CheckRedirect: http.DefaultClient.CheckRedirect,
			Jar:           http.DefaultClient.Jar,
			Timeout:       10 * time.Second,
		},
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

type urlOption string

var _ Option = urlOption("")

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

type timeoutOption time.Duration

var _ Option = timeoutOption(0)

func (t timeoutOption) apply(opts *options) error {
	if t < 0 {
		return fmt.Errorf("timeout is %d, needs to be > 0", t)
	}
	opts.HTTPClient.Timeout = time.Duration(t)
	return nil
}

// WithTimeout sets the timeout for signalling requests (/start, ...).
//
// The default is 10s.
//
// Note that you can also supply a timeout via [WithHTTPClient].
// Whichever option is provided last will take precedence.
func WithTimeout(t time.Duration) Option {
	return timeoutOption(t)
}

type httpClientOption struct {
	client *http.Client
}

var _ Option = httpClientOption{}

func (h httpClientOption) apply(opts *options) error {
	if h.client == nil {
		return errors.New("HTTP client must be non-nil")
	}
	opts.HTTPClient = h.client
	return nil
}

// WithHTTPClient sets a custom HTTP client to be used during signalling requests.
//
// Default is [http.DefaultClient].
//
// Note that providing [WithTimeout] after [WithHTTPClient] will overwrite this client's timeout value.
func WithHTTPClient(client *http.Client) Option {
	return httpClientOption{client: client}
}
