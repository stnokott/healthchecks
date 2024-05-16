package healthchecks

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func request(ctx context.Context, opts *options, body io.Reader, path ...string) error {
	fullPath := opts.RootURL.JoinPath(path...)

	req, err := newRequest(ctx, fullPath.String(), body)
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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		respBody = []byte("no information")
	}
	respStr := string(respBody)
	// if not 200, we return an error with the response body
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP response status %d: %s", resp.StatusCode, respStr)
	}
	// if the body is "OK", it's an actual good response
	if respStr != "OK" {
		return fmt.Errorf("HTTP response not OK: '%s'", respStr)
	}
	return nil
}

func newRequest(ctx context.Context, path string, body io.Reader) (*http.Request, error) {
	var method string
	if body == nil {
		method = "GET"
	} else {
		method = "POST"
	}

	return http.NewRequestWithContext(ctx, method, path, body)
}
