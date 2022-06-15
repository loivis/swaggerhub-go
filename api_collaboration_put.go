package swaggerhub

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *Client) APICollaborationPut(param APICollaborationPutParam) error {
	u := fmt.Sprintf("%s/apis/%s/%s/.collaboration", c.baseURL, param.Owner, param.API)
	log.Printf("request: %s -> %s", http.MethodPut, u)
	log.Printf("collaboration: %s", param.Body)
	resp, err := c.do(http.MethodPut, u, ContentType{Request: contentTypeJSON, Response: contentTypeJSON}, bytes.NewReader(param.Body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}

		log.Println(string(b))

		return fmt.Errorf("unexpected status code %d, want %d", code, http.StatusOK)
	}

	return nil
}

type APICollaborationPutParam struct {
	Owner string
	API   string
	Body  []byte
}
