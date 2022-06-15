package swaggerhub

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *Client) APIGet(param APIGetParam) ([]byte, error) {
	if !param.Resolved && param.Flatten {
		return nil, fmt.Errorf("use flatten only when resolved set true")
	}

	u := fmt.Sprintf("%s/apis/%s/%s/%s?resolved=%t&flatten=%t", c.baseURL, param.Owner, param.API, param.Version, param.Resolved, param.Flatten)
	log.Printf("request: %s -> %s", http.MethodGet, u)
	resp, err := c.do(http.MethodGet, u, ContentType{Response: param.ContentType.Response}, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d, want %d", code, http.StatusOK)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type APIGetParam struct {
	Owner   string
	API     string
	Version string
	// Set to true to get the resolved version of the API definition.
	// The content of all external $refs will be included in the resulting file.
	Resolved bool
	// Used only if resolved=true.
	// Flattening replaces all complex inline schemas with named entries in the components/schemas or definitions section.
	Flatten     bool
	ContentType ContentType
}
