package main

import (
	"io"
	"net/http"
)

type simpleHttpClient struct {
	httpClient *http.Client
}

func NewSimpleHttpClient(httpClient *http.Client) *simpleHttpClient {
	return &simpleHttpClient{httpClient: httpClient}
}

func (c *simpleHttpClient) GetContentFromUrl(url string) ([]byte, error) {
	resp, err := c.httpClient.Get(normalizeUrl(url))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}
