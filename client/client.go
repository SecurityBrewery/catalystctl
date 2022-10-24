package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SecurityBrewery/catalystctl/config"
	"io"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type CatalystClient struct {
	URL        *url.URL
	Token      string
	HTTPClient *http.Client
}

func New(url *url.URL, token string) (*CatalystClient, error) {
	u, err := config.URL(url)
	if err != nil {
		return nil, err
	}
	t, err := config.Token(token)
	if err != nil {
		return nil, err
	}

	return &CatalystClient{URL: u, Token: t, HTTPClient: http.DefaultClient}, nil
}

func (c *CatalystClient) Do(req *http.Request) (*http.Response, error) {
	req.URL = c.URL.ResolveReference(req.URL)

	if req.Header == nil {
		req.Header = make(http.Header)
	}
	req.Header.Set("PRIVATE-TOKEN", c.Token)
	return c.HTTPClient.Do(req)
}

func (c *CatalystClient) Get(path string) (*http.Response, error) {
	return c.Do(&http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Path: path},
	})
}

func (c *CatalystClient) PostJSON(path string, header http.Header, data any) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/json")
	return c.Post(path, header, bytes.NewReader(body))
}

func (c *CatalystClient) Post(path string, header http.Header, body io.Reader) (*http.Response, error) {
	return c.Do(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Path: path},
		Header: header,
		Body:   io.NopCloser(body),
	})
}

func (c *CatalystClient) Version() (string, error) {
	resp, err := c.Get("/api/settings")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", errors.New("invalid token")
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(responseData))
	}

	catalystVersion := gjson.GetBytes(responseData, "version").String()
	return catalystVersion, nil
}
