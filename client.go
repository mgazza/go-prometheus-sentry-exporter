package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	DefaultClient HttpClient = http.DefaultClient
)

type Client struct {
	baseUrl string
	token   string
	org     string
}

func NewSentryClient(baseUrl string, org string, token string) *Client {
	return &Client{
		baseUrl: baseUrl,
		org:     org,
		token:   token,
	}
}

type unexpectedResponseError struct {
	code int
}

func (r unexpectedResponseError) Error() string {
	return fmt.Sprintf("unexpected response (%d)", r.code)
}

func (c *Client) GetProjects() (*[]ProjectResp, error) {
	uri, err := url.Parse(c.baseUrl)
	if err != nil {
		return nil, err
	}
	uri.Path = fmt.Sprintf("%s/projects/", uri.Path)
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		return nil, err
	}
	req.URL = uri
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/json")

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, unexpectedResponseError{
			code: resp.StatusCode,
		}
	}

	decoder := json.NewDecoder(resp.Body)

	var results []ProjectResp
	err = decoder.Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

func (c *Client) GetIssues(projectSlug string) (*[]IssueResp, error) {
	uri, err := url.Parse(c.baseUrl)
	if err != nil {
		return nil, err
	}
	uri.Path = fmt.Sprintf("%s/projects/%s/%s/issues/", uri.Path, c.org, projectSlug)
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		return nil, err
	}
	req.URL = uri
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/json")

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, unexpectedResponseError{
			code: resp.StatusCode,
		}
	}

	decoder := json.NewDecoder(resp.Body)

	var results []IssueResp
	err = decoder.Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}
