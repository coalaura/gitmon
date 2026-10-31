package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type GitHubResponse struct {
	Items []GitHubItem `json:"items"`
}

type GitHubItem struct {
	User      GitHubUser `json:"user"`
	HtmlURL   string     `json:"html_url"`
	Title     string     `json:"title"`
	State     string     `json:"state"`
	Number    int64      `json:"number"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Draft     bool       `json:"draft"`
}

type GitHubUser struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HtmlURL   string `json:"html_url"`
}

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func (t Target) Fetch() ([]GitHubItem, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/search/issues", nil)
	if err != nil {
		return nil, err
	}

	query := make(url.Values)

	query.Set("q", t.GetFilter())
	query.Set("per_page", "100")

	if t.Sort == "" {
		query.Set("sort", "created_at")
	} else {
		query.Set("sort", t.Sort)
	}

	if t.Order == "" {
		query.Set("order", "desc")
	} else {
		query.Set("order", t.Order)
	}

	req.URL.RawQuery = query.Encode()

	if config.Server.GitHubToken != "" {
		req.Header.Set("Authentication", "Bearer "+config.Server.GitHubToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	var data GitHubResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Items, nil
}

func (t Target) GetFilter() string {
	filter := strings.TrimSpace(t.Filter)

	switch t.Type {
	case "issues":
		filter += " is:issue"
	case "pull-requests":
		filter += " is:pr"
	}

	return filter
}

func RunFetchLoop() {

}
