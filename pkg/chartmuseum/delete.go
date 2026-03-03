package chartmuseum

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type deleteResponse struct {
	Deleted bool   `json:"deleted"`
	Error   string `json:"error"`
}

// DeleteChartVersion sends a DELETE request to remove a specific chart version from ChartMuseum.
func (c *Client) DeleteChartVersion(name, version string) error {
	u, err := url.Parse(c.opts.url)
	if err != nil {
		return fmt.Errorf("invalid URL %q: %w", c.opts.url, err)
	}

	contextPath := strings.TrimSuffix(c.opts.contextPath, "/")
	u.Path, err = url.JoinPath(u.Path, contextPath, "api", "charts", url.PathEscape(name), url.PathEscape(version))
	if err != nil {
		return fmt.Errorf("failed to join url path: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	setAuth(req, c.opts)

	resp, err := c.Client.Do(req) //nolint:staticcheck
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var dr deleteResponse
		if err := json.Unmarshal(body, &dr); err != nil {
			return fmt.Errorf("unexpected response from ChartMuseum: %s", string(body))
		}
		if !dr.Deleted {
			msg := dr.Error
			if msg == "" {
				msg = "chart not found or already deleted"
			}
			return fmt.Errorf("%s-%s: %s", name, version, msg)
		}
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("%s-%s: not found in repository", name, version)
	default:
		return chartmuseumError(body, resp.StatusCode)
	}
}

func setAuth(req *http.Request, opts options) {
	if opts.authHeader != "" && opts.accessToken != "" {
		req.Header.Set(opts.authHeader, opts.accessToken)
	} else if opts.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+opts.accessToken)
	} else if opts.username != "" || opts.password != "" {
		req.SetBasicAuth(opts.username, opts.password)
	}
}

func chartmuseumError(body []byte, statusCode int) error {
	var dr deleteResponse
	if err := json.Unmarshal(body, &dr); err == nil && dr.Error != "" {
		return fmt.Errorf("ChartMuseum %d: %s", statusCode, dr.Error)
	}
	// Fallback to raw body if it's not JSON or has no error field.
	msg := strings.TrimSpace(string(body))
	if msg == "" {
		msg = http.StatusText(statusCode)
	}
	return fmt.Errorf("ChartMuseum %d: %s", statusCode, msg)
}
