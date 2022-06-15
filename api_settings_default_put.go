package swaggerhub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *Client) APISettingsDefaultPut(param APISettingsDefaultPutParam) error {
	u := fmt.Sprintf("%s/apis/%s/%s/settings/default", c.baseURL, param.Owner, param.API)
	log.Printf("request: %s -> %s, version %q", http.MethodPut, u, param.Version)
	v := version{param.Version}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	log.Println("request body:", string(b))
	resp, err := c.do(http.MethodPut, u, ContentType{Request: contentTypeJSON}, bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		log.Printf("%s/%s default version set to %s", param.Owner, param.API, param.Version)
		return nil
	default:
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("error reading response body:", err)
		}
		log.Println("unexpected response:", string(b))
		return fmt.Errorf("unexpected status code %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

type APISettingsDefaultPutParam struct {
	Owner   string
	API     string
	Version string
}
