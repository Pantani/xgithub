package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v50/github"

	"github.com/ignite/bounty/xgithub"
)

const (
	accessToken = "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"
	igniteTopic = "ignite-plugin"
)

func bountyPOC() {
	var (
		// client context
		ctx = context.Background()
		c   = xgithub.NewClient(ctx, accessToken)

		// create a context with timeout for queries
		queryCtx, cancel = context.WithTimeout(ctx, 5*time.Second)

		// sort the repo by the stars
		opts = &github.SearchOptions{Sort: "stars", Order: "desc"}
	)
	defer cancel()

	result, err := c.RepoQuery(queryCtx, opts,
		// We should define a static topic name to always search like "ignite-plugin" or "ignite-cli-plugin".
		// https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories#search-by-topic
		xgithub.NewQuery("topic", igniteTopic),
		// We should only list a repo with a few stars to have the basic criteria.
		// https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories#search-by-number-of-stars
		xgithub.NewQuery("stars", ">5"),
		// It only shows the Golang repositories to run non-Golang applications.
		// https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories#search-by-language
		xgithub.NewQuery("language", "go"),
		// We should define a license to list the user plugin in the CLI. We can also automatically scaffold the
		// license with a new plugin template.
		// https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories#search-by-license
		xgithub.NewQuery("license", "MIT"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", result)
}
