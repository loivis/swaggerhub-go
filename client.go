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

func (c *Client) newRequest(method, u string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.apiKey)

	return req, nil
}

func (c *Client) do(method, u, format string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.apiKey)
	accept := contentTypeJSON
	if format != "" {
		accept = "application/" + format
	}
	req.Header.Add("Accept", accept)

	return c.hc.Do(req)
}

const (
	contentTypeJSON = "application/json"
)

// type APIList struct {
// 	Name        string
// 	Description string
// 	URL         string
// 	Offset      int
// 	TotalCount  int
// 	APIs        []API
// }

// type API struct {
// 	Name        string
// 	Description string
// 	Tags        []string
// 	Properties  []APIProperty
// }

// type APIProperty map[string]string

type version struct {
	Version string `json:"version"`
}

type published struct {
	Published bool `json:"published"`
}

// type PublishSettings struct {
// 	Published bool
// }
