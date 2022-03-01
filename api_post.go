package swaggerhub

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *Client) APIPost(param APIPostParam) error {
	u := fmt.Sprintf("%s/apis/%s/%s?isPrivate=%t", c.baseURL, param.Owner, param.API, param.IsPrivate)
	if param.Version != "" {
		u += "&version=" + param.Version
	}
	log.Printf("request: %s -> %s", http.MethodPost, u)
	req, err := c.newRequest(http.MethodPost, u, param.Body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", param.RequestContentType)

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		log.Printf("%s/%s updated", param.Owner, param.API)
		return nil
	case http.StatusCreated:
		log.Printf("%s/%s created", param.Owner, param.API)
		return nil
	default:
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("error reading response body:", err)
		}
		log.Println("unexpected response:", string(b))
		return fmt.Errorf("unexpected status code %d, want %d or %d", resp.StatusCode, http.StatusOK, http.StatusCreated)
	}
}

type APIPostParam struct {
	Owner              string
	API                string
	Version            string
	IsPrivate          bool
	RequestContentType string
	Body               io.Reader
}
