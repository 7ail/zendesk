package zendesk

import (
	"errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type accountResponse struct {
	Url          string `json:"url"`
	Name         string `json:"name"`
	Sandbox      bool   `json:"sandbox"`
	Subdomain    string `json:"subdomain"`
	TimeFormat   int64  `json:"time_format"`
	TimeZone     string `json:"time_zone"`
	OwnerID      int64  `json:"owner_id"`
	Multiproduct bool   `json:"multiproduct"`
}

func (g *gateway) Url() (string, error) {
	resp, err := g.account()
	if err != nil {
		return "", fmt.Errorf("g.account(): %w", err)
	}

	if resp.Url == "" {
		return "", errors.Join(ErrUnexpected, fmt.Errorf("g.account(): %w", ErrNotFound))
	}

	return resp.Url, nil
}

func (g *gateway) Subdomain() (string, error) {
	resp, err := g.account()
	if err != nil {
		return "", fmt.Errorf("g.account(): %w", err)
	}

	if resp.Subdomain == "" {
		return "", errors.Join(ErrUnexpected, fmt.Errorf("g.account(): %w", ErrNotFound))
	}

	return resp.Subdomain, nil
}

func (g *gateway) account() (accountResponse, error) {
	result := accountResponse{}

	requestUrl := fmt.Sprintf("http://%s.%s.com/api/v2/account.json", g.subdomain, g.host)
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return result, errors.Join(ErrUnexpected, fmt.Errorf("http.NewRequest: %w", err))
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return result, errors.Join(ErrUnexpected, fmt.Errorf("g.client.Do: %w", err))
	}

	if resp.StatusCode != http.StatusOK {
		return result, convertToError(resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, errors.Join(ErrUnexpected, fmt.Errorf("ioutil.ReadAll: %w", err))
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, errors.Join(ErrUnexpected, fmt.Errorf("json.Unmarshal: %w", err))
	}

	return result, nil
}
