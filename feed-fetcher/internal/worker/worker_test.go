package worker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/config"
	"github.com/vladislav-chunikhin/feed-fetcher/internal/worker/mocks"
)

func TestWorker_Run(t *testing.T) {
	ctx := context.TODO()
	logger := mocks.NewMockLogger()

	feedFetcher := mocks.NewFeedFetcher(t)
	feedFetcher.EXPECT().
		Fetch(ctx).
		Return(nil)

	worker, err := NewWorker(feedFetcher, config.WorkerConfig{Interval: 1 * time.Second}, logger)
	require.NoError(t, err)

	err = worker.Run(ctx)
	require.NoError(t, err)

	// Waiting for the cron task to finish
	time.Sleep(2 * time.Second)

	worker.Stop()
}
