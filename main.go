package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v50/github"

	"github.com/ignite/bounty/xgithub"
)

func main() {
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
		//xgithub.NewQuery("q", "// this line is used by starport scaffolding #"),
		xgithub.NewQuery("'// this line is used by starport scaffolding #'+in", "file"),
		xgithub.NewQuery("language", "go"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", result)
}
