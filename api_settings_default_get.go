package swaggerhub

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) APISettingsDefaultGet(param *APISettingsDefaultGetParam) (string, error) {
	u := fmt.Sprintf("%s/apis/%s/%s/settings/default", c.baseURL, param.Owner, param.API)
	resp, err := c.do(http.MethodGet, u, ContentType{Response: contentTypeJSON}, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d, want %d", code, http.StatusOK)
	}

	var ver version
	if err := json.NewDecoder(resp.Body).Decode(&ver); err != nil {
		return "", err
	}

	return ver.Version, nil
}

type APISettingsDefaultGetParam struct {
	Owner string
	API   string
}
