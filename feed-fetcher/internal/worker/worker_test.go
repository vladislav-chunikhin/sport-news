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
	feedFetcher := mocks.NewFeedFetcher()

	worker, err := NewWorker(feedFetcher, config.WorkerConfig{Interval: 1 * time.Second}, logger)
	require.NoError(t, err)

	err = worker.Run(ctx)
	require.NoError(t, err)
}

func TestWorker_Run_Invalid_Interval(t *testing.T) {
	logger := mocks.NewMockLogger()
	feedFetcher := mocks.NewFeedFetcher()

	_, err := NewWorker(feedFetcher, config.WorkerConfig{Interval: 1 * time.Microsecond}, logger)
	require.Error(t, err)
	require.Equal(t, "the interval must consist of whole seconds, for example, 1s, 2s, 3s, etc. Invalid value is 50ms, as it corresponds to 0.05s", err.Error())
}

func TestWorker_NewWorker_Nil_Fetcher(t *testing.T) {
	logger := mocks.NewMockLogger()
	_, err := NewWorker(nil, config.WorkerConfig{Interval: 1 * time.Microsecond}, logger)
	require.Error(t, err)
	require.Equal(t, "nil fetcher", err.Error())
}
