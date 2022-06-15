package swaggerhub

import (
	"io"
	"net/http"
)

type Client struct {
	baseURL string
	apiKey  string
	hc      *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		hc:      http.DefaultClient,
	}
}

func (c *Client) newRequest(method, u string, ct ContentType, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.apiKey)

	if ct.Request != "" {
		req.Header.Set("Content-Type", ct.Request)
	}

	if ct.Response != "" {
		req.Header.Set("Accept", ct.Response)
	}

	return req, nil
}

func (c *Client) do(method, u string, ct ContentType, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(method, u, ct, body)
	if err != nil {
		return nil, err
	}

	return c.hc.Do(req)
}

const (
	contentTypeJSON = "application/json"
)

type version struct {
	Version string `json:"version"`
}

type published struct {
	Published bool `json:"published"`
}

type ContentType struct {
	Request  string
	Response string
}
