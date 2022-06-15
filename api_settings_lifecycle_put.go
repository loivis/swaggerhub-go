package swaggerhub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *Client) APISettingsLifecyclePut(param APISettingsLifecyclePutParam) error {
	u := fmt.Sprintf("%s/apis/%s/%s/%s/settings/lifecycle?force=%t", c.baseURL, param.Owner, param.API, param.Version, param.Force)
	log.Printf("request: %s -> %s", http.MethodPut, u)
	v := published{param.Published}
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
		log.Printf("%s/%s/%s published: %t", param.Owner, param.API, param.Version, param.Published)
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

type APISettingsLifecyclePutParam struct {
	Owner     string
	API       string
	Version   string
	Published bool
	Force     bool
}
