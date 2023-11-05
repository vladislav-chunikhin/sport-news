package app

import (
	"context"
	"fmt"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/fetcher"
	"github.com/vladislav-chunikhin/feed-fetcher/internal/provider/htafc"
)

func (c *Container) GetFeedProvider(_ context.Context) fetcher.FeedProvider {
	if c.feedProvider == nil {
		var err error
		c.feedProvider, err = htafc.New(&c.Config().FeedProviders.Htafc, c.Logger())
		if err != nil {
			panic(fmt.Errorf("failed to create feed provider client: %v", err))
		}
	}

	return c.feedProvider
}
