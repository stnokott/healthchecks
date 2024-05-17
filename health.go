// Package healthchecks is a wrapper around the healthchecks.io endpoints.
package healthchecks

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Project organizes multiple checks in a common project.
//
// Use [NewProject] for obtaining a new instance.
type Project struct {
	pingKey string
	opts    *options
}

// NewProject creates a new instance of [Project].
//
// Sending signals via a project requires the project's "ping key" and the check's slug.
// The ping key can be created under your project's settings.
func NewProject(pingKey string, opts ...Option) (*Project, error) {
	if pingKey == "" {
		return nil, errors.New("project ping key must not be empty")
	}

	options, err := optsFromDefaults(opts)
	if err != nil {
		return nil, err
	}

	return &Project{
		pingKey: pingKey,
		opts:    options,
	}, nil
}

// Notifier sends signals to a check.
type Notifier interface {
	// Start sends the "start" signal to the check.
	Start(ctx context.Context) error
	// Success sends the "success" signal to the check.
	Success(ctx context.Context) error
	// Fail sends the "fail" signal to the check.
	Fail(ctx context.Context) error
	// Log sends the "log" signal with the attached message to the check.
	Log(ctx context.Context, msg string) error
	// ExitStatus sends the "exit-status" signal with the exit code to the check.
	//
	// Success or failure of the check is determined by the exit code.
	ExitStatus(ctx context.Context, code int) error
}

// TODO: exit status https://healthchecks.io/docs/http_api/#exitcode-uuid

// Start sends the "start" signal to the project's check identified by slug.
func (p *Project) Start(ctx context.Context, slug string) error {
	return request(ctx, p.opts, nil, p.pingKey, slug, "/start")
}

// Success sends the "success" signal to the project's check identified by slug.
func (p *Project) Success(ctx context.Context, slug string) error {
	return request(ctx, p.opts, nil, p.pingKey, slug)
}

// Fail sends the "fail" signal to the project's check identified by slug.
func (p *Project) Fail(ctx context.Context, slug string) error {
	return request(ctx, p.opts, nil, p.pingKey, slug, "/fail")
}

// Log sends the "log" signal with the attached message to the project's check identified by slug.
func (p *Project) Log(ctx context.Context, slug string, msg string) error {
	return request(ctx, p.opts, strings.NewReader(msg), p.pingKey, slug, "/log")
}

// ExitStatus sends the "exit-status" signal with the exit code to the project's check identified by slug.
//
// Success or failure of the check is determined by the exit code.
func (p *Project) ExitStatus(ctx context.Context, slug string, code int) error {
	return request(ctx, p.opts, nil, p.pingKey, slug, "/", strconv.Itoa(code))
}

// Slug creates a new [Notifier] for a check in this [Project], indentified by its slug.
func (p *Project) Slug(slug string) Notifier {
	return &Check{
		path: p.pingKey + "/" + slug,
		opts: p.opts,
	}
}

// compile-time interface implementation check
var _ Notifier = (*Check)(nil)

// Check is an individual check, either obtained via [NewUUID] or [Project.Slug] (via [NewProject]).
//
// It implements [Notifier].
type Check struct {
	path string
	opts *options
}

// NewUUID creates a new instance of [Check], identified by its UUID.
//
// The UUID is in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func NewUUID(uuid string, opts ...Option) (*Check, error) {
	if uuid == "" {
		return nil, errors.New("uuid must not be empty")
	}
	options, err := optsFromDefaults(opts)
	if err != nil {
		return nil, err
	}

	return &Check{
		path: "/" + uuid,
		opts: options,
	}, nil
}

// FromURL constructs a new Notifier from the URL of a single (UUID-based) check.
//
// The URL should be in the format http(s)://example.com/uuid.
//
// Note that any option which overrides the URL will be ignored.
func FromURL(u string, opts ...Option) (Notifier, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}

	// apply options
	options, err := optsFromDefaults(opts)
	if err != nil {
		return nil, err
	}

	// override URL
	options.RootURL, _ = url.Parse(parsed.Scheme + "://" + parsed.Host) //nolint:errcheck

	return &Check{
		path: parsed.Path,
		opts: options,
	}, nil
}

// Start sends the "start" signal to the check identified by its uuid.
func (c *Check) Start(ctx context.Context) error {
	return request(ctx, c.opts, nil, c.path, "/start")
}

// Success sends the "success" signal to the check identified by its uuid.
func (c *Check) Success(ctx context.Context) error {
	return request(ctx, c.opts, nil, c.path)
}

// Fail sends the "fail" signal to the check identified by its uuid.
func (c *Check) Fail(ctx context.Context) error {
	return request(ctx, c.opts, nil, c.path, "/fail")
}

// Log sends the "log" signal with the attached message to the check identified by its uuid.
func (c *Check) Log(ctx context.Context, msg string) error {
	return request(ctx, c.opts, strings.NewReader(msg), c.path, "/log")
}

// ExitStatus sends the "exit-status" signal with the exit code to the check identified by its uuid.
//
// Success or failure of the check is determined by the exit code.
func (c *Check) ExitStatus(ctx context.Context, code int) error {
	return request(ctx, c.opts, nil, c.path, "/", strconv.Itoa(code))
}
