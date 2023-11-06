package mocks

import "context"

type FeedFetcher struct{}

func NewFeedFetcher() *FeedFetcher {
	return &FeedFetcher{}
}

func (m *FeedFetcher) Fetch(_ context.Context) error {
	return nil
}
