package healthchecks

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func request(ctx context.Context, opts *options, path ...string) error {
	fullPath := opts.RootURL.JoinPath(path...)

	req, err := http.NewRequestWithContext(ctx, "GET", fullPath.String(), nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	// body is required to differentiate between 200s (OK, not found, rate limited etc.).
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		body = []byte("no information")
	}
	bodyStr := string(body)
	// if not 200, we return an error with the response body
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP response status %d: %s", resp.StatusCode, bodyStr)
	}
	// if the body is "OK", it's an actual good response
	if bodyStr != "OK" {
		return fmt.Errorf("HTTP response not OK: '%s'", bodyStr)
	}
	return nil
}
