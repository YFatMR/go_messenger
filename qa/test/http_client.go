package test

import (
	"io"
	"net/http"
)

type HttpClient struct {
	http.Client
	BaseUrl string
}

func (c *HttpClient) Get(url string) (*http.Response, error) {
	return c.Client.Get(c.BaseUrl + url)
}

func (c *HttpClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	return c.Client.Post(c.BaseUrl+url, contentType, body)
}

func (c *HttpClient) Patch(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPatch, c.BaseUrl+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}
