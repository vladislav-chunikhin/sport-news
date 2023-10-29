package fetcher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vladislav-chunikhin/lib-go/pkg/logger"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/mapper"
	"github.com/vladislav-chunikhin/feed-fetcher/internal/provider/htafc"
)

type Producer interface {
	PublishHtafcFeed(ctx context.Context, message []byte) error
}

type FeedProvider interface {
	GetLatestNews(ctx context.Context) (*htafc.NewListInformation, error)
}

type Fetcher struct {
	producer     Producer
	feedProvider FeedProvider
	logger       logger.Logger
}

func NewFetcher(producer Producer, htafcFeedProvider FeedProvider, logger logger.Logger) (*Fetcher, error) {
	if producer == nil {
		return nil, fmt.Errorf("nil producer")
	}

	if htafcFeedProvider == nil {
		return nil, fmt.Errorf("nil htafc feed provider")
	}

	return &Fetcher{producer: producer, feedProvider: htafcFeedProvider, logger: logger}, nil
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	latestNews, err := f.feedProvider.GetLatestNews(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest news: %w", err)
	}

	if latestNews == nil {
		return fmt.Errorf("nil result from feed client")
	}

	if len(latestNews.NewsletterNews.NewsletterNews) == 0 {
		return fmt.Errorf("no feeds from provider")
	}

	articleItems := mapper.ToArticleItems(latestNews, f.logger)
	if articleItems == nil {
		return nil
	}

	var msg []byte
	msg, err = json.Marshal(articleItems)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	if err = f.producer.PublishHtafcFeed(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
