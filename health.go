// Package healthchecks is a wrapper around the healthchecks.io endpoints.
package healthchecks

import (
	"context"
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
	options, err := fromDefaults(opts)
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
}

// TODO: log https://healthchecks.io/docs/http_api/#log-uuid
// TODO: exit status https://healthchecks.io/docs/http_api/#exitcode-uuid

// Start sends the "start" signal to the project's check identified by slug.
func (p *Project) Start(ctx context.Context, slug string) error {
	return request(ctx, p.opts, p.pingKey, slug, "/start")
}

// Success sends the "success" signal to the project's check identified by slug.
func (p *Project) Success(ctx context.Context, slug string) error {
	return request(ctx, p.opts, p.pingKey, slug)
}

// Fail sends the "fail" signal to the project's check identified by slug.
func (p *Project) Fail(ctx context.Context, slug string) error {
	return request(ctx, p.opts, p.pingKey, slug, "/fail")
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

// Check is an individual check, either obtained via [NewUUID].
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
	options, err := fromDefaults(opts)
	if err != nil {
		return nil, err
	}

	return &Check{
		path: uuid,
		opts: options,
	}, nil
}

// Start sends the "start" signal to the check identified by its uuid.
func (c *Check) Start(ctx context.Context) error {
	return request(ctx, c.opts, c.path, "/start")
}

// Success sends the "success" signal to the check identified by its uuid.
func (c *Check) Success(ctx context.Context) error {
	return request(ctx, c.opts, c.path)
}

// Fail sends the "fail" signal to the check identified by its uuid.
func (c *Check) Fail(ctx context.Context) error {
	return request(ctx, c.opts, c.path, "/fail")
}
