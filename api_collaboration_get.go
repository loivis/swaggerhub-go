package swaggerhub

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *Client) APICollaborationGet(param APICollaborationGetParam) ([]byte, error) {
	// TODO: bug?
	// hardcode expandTeams to be true.
	// setting it to false or missing the parameter results in 400.
	u := fmt.Sprintf("%s/apis/%s/%s/.collaboration?expandTeams=true", c.baseURL, param.Owner, param.API)
	log.Printf("request: %s -> %s", http.MethodGet, u)
	resp, err := c.do(http.MethodGet, u, ContentType{Response: contentTypeJSON}, nil)
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

type APICollaborationGetParam struct {
	Owner string
	API   string
}
