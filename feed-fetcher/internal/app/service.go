package app

import (
	"context"
	"fmt"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/fetcher"
)

func (c *Container) GetFetcher(ctx context.Context) *fetcher.Fetcher {
	if c.feedFetcher == nil {
		var err error
		c.feedFetcher, err = fetcher.NewFetcher(c.Producer(), c.GetFeedProvider(ctx), c.Logger())
		if err != nil {
			panic(fmt.Errorf("failed to create fetcher: %v", err))
		}
	}

	return c.feedFetcher
}
