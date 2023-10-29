package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/jasonlvhit/gocron"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/config"
	"github.com/vladislav-chunikhin/lib-go/pkg/logger"
)

const defaultInterval = 30

type FeedFetcher interface {
	Fetch(ctx context.Context) error
}

type Worker struct {
	fetcher FeedFetcher
	job     *gocron.Job
	logger  logger.Logger
}

func NewWorker(fetcher FeedFetcher, config config.WorkerConfig, logger logger.Logger) (*Worker, error) {
	if fetcher == nil {
		return nil, fmt.Errorf("nil fetcher")
	}

	if config.Interval == 0 {
		logger.Warnf("interval equals zero, will be used default value: %ds", defaultInterval)
		config.Interval = defaultInterval * time.Second
	}

	seconds := config.Interval.Seconds()

	secondsAsUint64 := uint64(seconds)

	if float64(secondsAsUint64) != seconds {
		return nil, fmt.Errorf("the interval must consist of whole seconds, for example, 1s, 2s, 3s, etc. Invalid value is 50ms, as it corresponds to 0.05s")
	}

	job := gocron.Every(secondsAsUint64).Seconds()

	return &Worker{
		fetcher: fetcher,
		job:     job,
		logger:  logger,
	}, nil
}

func (w *Worker) Run(ctx context.Context) error {
	if w.job == nil {
		return fmt.Errorf("nil job")
	}

	if err := w.job.Do(func() {
		w.logger.Debugf("starting to fetch data from the provider...")
		if err := w.fetcher.Fetch(ctx); err != nil {
			w.logger.Errorf("failed to fetch feeds: %v", err)
			return
		}
	}); err != nil {
		return err
	}

	gocron.Start()
	w.logger.Debugf("starting the scheduled job...")

	return nil
}

func (w *Worker) Stop() {
	if w.job == nil {
		return
	}

	gocron.Clear()
	w.logger.Debugf("deleting the scheduled job...")
}
