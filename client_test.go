package main_test

import (
	"bytes"
	ex "github.com/mgazza/go-sentry-prometheus-exporter"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type ClientMock struct {
	Url          string
	Headers      map[string]string
	T            *testing.T
	ResponseBody string
	Code         int
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	assert.Equal(c.T, c.Url, req.URL.String())

	for k, v := range c.Headers {
		assert.Equal(c.T, v, req.Header.Get(k))
	}

	responseBody := ioutil.NopCloser(bytes.NewReader([]byte(c.ResponseBody)))

	result := &http.Response{
		StatusCode: c.Code,
		Body:       responseBody,
	}
	return result, nil
}

func TestClient_GetProjects(t *testing.T) {
	type fields struct {
		baseUrl string
		token   string
		org     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *[]ex.ProjectResp
		mock    *ClientMock
		wantErr bool
	}{
		{
			name: "When getting projects",
			fields: fields{
				baseUrl: "https://sentry.io/api/0",
				token:   "thetoken",
				org:     "n/a",
			},
			mock: &ClientMock{
				Code: 200,
				Headers: map[string]string{
					"Authorization": "Bearer thetoken",
				},
				Url:          "https://sentry.io/api/0/projects/",
				ResponseBody: `[{"slug":"theslug","name":"thename"}]`,
			},
			want: &[]ex.ProjectResp{
				{
					Slug: "theslug",
					Name: "thename",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ex.NewSentryClient(tt.fields.baseUrl, tt.fields.org, tt.fields.token)

			oldClient := ex.DefaultClient
			defer func() { ex.DefaultClient = oldClient }()
			ex.DefaultClient = tt.mock
			tt.mock.T = t

			got, err := c.GetProjects()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProjects() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetIssues(t *testing.T) {
	type fields struct {
		baseUrl string
		token   string
		org     string
	}
	type args struct {
		projectSlug string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    *ClientMock
		want    *[]ex.IssueResp
		wantErr bool
	}{
		{
			name: "When getting issues",
			fields: fields{
				baseUrl: "https://sentry.io/api/0",
				token:   "thetoken",
				org:     "theorg",
			},
			args: args{
				projectSlug: "theproject",
			},
			mock: &ClientMock{
				Code: 200,
				Headers: map[string]string{
					"Authorization": "Bearer thetoken",
				},
				Url:          "https://sentry.io/api/0/projects/theorg/theproject/issues/",
				ResponseBody: `[{"type":"Issue","permalink":"https://...","level":"WARN","count":"1"}]`,
			},
			want: &[]ex.IssueResp{
				{
					Logger:    "",
					Type:      "Issue",
					Permalink: "https://...",
					Level:     "WARN",
					Count:     "1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ex.NewSentryClient(tt.fields.baseUrl, tt.fields.org, tt.fields.token)

			oldClient := ex.DefaultClient
			defer func() { ex.DefaultClient = oldClient }()
			ex.DefaultClient = tt.mock
			tt.mock.T = t

			got, err := c.GetIssues(tt.args.projectSlug)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIssues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIssues() got = %v, want %v", got, tt.want)
			}
		})
	}
}
