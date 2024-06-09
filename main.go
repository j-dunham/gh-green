package main

import (
	"fmt"
	"github.com/cli/go-gh/v2/pkg/api"
	"log"
	"time"
)

// Need for wrapping the from argument so the type matches the expected input type
type DateTime struct{ time.Time }

func ExampleContributionGraphQL() {
	opts := api.ClientOptions{
		EnableCache: true,
		Timeout:     5 * time.Second,
	}
	client, err := api.NewGraphQLClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	var Query struct {
		Viewer struct {
			ContributionsCollection struct {
				TotalCommitContributions int
				HasAnyContributions      bool
			} `graphql:"contributionsCollection(from: $from)"`
		} `graphql:"viewer"`
	}

	err = client.Query("contributionQuery", &Query, map[string]interface{}{"from": DateTime{time.Now()}})
	if err != nil {
		log.Fatal(err)
	}

	isGreen := Query.Viewer.ContributionsCollection.HasAnyContributions
	if isGreen {
		// print sparkle cake emoji
		fmt.Println("🍰 You are green!")
	} else {
		fmt.Println("You haven't made any commits today... yet!")
	}
}

func main() {
	ExampleContributionGraphQL()
}
