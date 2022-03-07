package main

type ProjectResp struct {
	// there are lots of other field that we do not care about
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type IssueResp struct {
	// there are lots of other field that we do not care about
	Logger    string `json:"logger"`
	Type      string `json:"type"`
	Permalink string `json:"permalink"`
	Level     string `json:"level"`
	Count     string `json:"count"`
}
