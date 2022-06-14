package swaggerhub

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *Client) DomainGet(param DomainGetParam) ([]byte, error) {
	u := fmt.Sprintf("%s/domains/%s/%s/%s", c.baseURL, param.Owner, param.Domain, param.Version)
	log.Printf("request: %s -> %s", http.MethodGet, u)
	resp, err := c.do(http.MethodGet, u, param.Format, nil)
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

type DomainGetParam struct {
	Owner   string
	Domain  string
	Version string
	Format  string
}
