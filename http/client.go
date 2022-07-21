package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chadwpetersen/ind/errors"
	"github.com/chadwpetersen/ind/log"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) GET(ctx context.Context, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.parseURL(url), nil)
	if err != nil {
		return nil, 0, err
	}

	return c.Do(req)
}

func (c *Client) POST(ctx context.Context, url string, body any) ([]byte, int, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.parseURL(url), bytes.NewBuffer(b))
	if err != nil {
		return nil, 0, err
	}

	return c.Do(req)
}

func (c *Client) Do(req *http.Request) ([]byte, int, error) {
	req.Header.Set("Content-Type", `application/json`)

	log.Debug("Performing HTTP request", log.WithLabels(
		map[string]any{
			"method": req.Method,
			"url":    req.URL.String(),
		},
	))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	// Clean-up the http response.
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	data = bytes.ReplaceAll(data, []byte(")]}',\n"), []byte(""))

	if resp.StatusCode != http.StatusOK {
		return data, resp.StatusCode, errors.Wrap(ErrInvalidHTTPStatus, map[string]any{
			"status": resp.StatusCode,
			"body":   string(data),
		})
	}

	return data, resp.StatusCode, nil
}

func (c *Client) parseURL(url string) string {
	return fmt.Sprintf("%v/%s", c.baseURL, url)
}
