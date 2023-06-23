package xgithub

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type (
	Client struct {
		ts     oauth2.TokenSource
		tc     *http.Client
		client *github.Client
	}

	Query struct {
		Qualifier string
		Value     string
	}

	RepoResult struct {
		ID          int64
		Name        string
		Description string
		GitURL      string
		HTMLURL     string
		// ...
	}
)

// RepoQuery runs a repository query search based on the parameters and returns the result or an error.
func (c *Client) RepoQuery(ctx context.Context, opts *github.SearchOptions, queries ...Query) ([]RepoResult, error) {
	query, err := CreateQuery(queries...)
	if err != nil {
		return nil, err
	}

	result, _, err := c.client.Search.Repositories(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	return parseRepoResult(result.Repositories), nil
}

// NewQuery creates a new GitHub query.
func NewQuery(qualifier, value string) Query {
	return Query{
		Qualifier: qualifier,
		Value:     value,
	}
}

// NewClient creates a new GitHub client.
func NewClient(ctx context.Context, accessToken string) *Client {
	var (
		ts     = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
		tc     = oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	)
	return &Client{
		ts:     ts,
		tc:     tc,
		client: client,
	}
}

// CreateQuery creates a new query string based on the query qualifiers, avoiding a repeat.
func CreateQuery(queries ...Query) (string, error) {
	var (
		qMap = make(map[string]string)
		q    = ""
	)
	for _, query := range queries {
		if _, ok := qMap[query.Qualifier]; ok {
			return "", fmt.Errorf("duplicated qualifier: %s", query.Qualifier)
		}
		qMap[query.Qualifier] = query.Value
		q += fmt.Sprintf(" %s:%s", query.Qualifier, query.Value)
	}
	return q, nil
}

// parseRepoResult parse the []*github.Repository into the local struct RepoResult.
func parseRepoResult(repos []*github.Repository) []RepoResult {
	result := make([]RepoResult, len(repos))
	for i, repo := range repos {
		result[i] = RepoResult{
			ID:      *repo.ID,
			Name:    *repo.Name,
			GitURL:  *repo.GitURL,
			HTMLURL: *repo.HTMLURL,
			// ...
		}
		if repo.Description != nil {
			result[i].Description = *repo.Description
		}
	}
	return result
}
