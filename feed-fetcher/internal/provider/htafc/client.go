package htafc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/vladislav-chunikhin/lib-go/pkg/httpclient"
	"github.com/vladislav-chunikhin/lib-go/pkg/logger"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/config"
)

const (
	basicAPIPath      = "/api/incrowd"
	getLatestNewsPath = "getnewlistinformation"

	providerName = "htafc"

	countParam       = "count"
	defaultBatchSize = 50
)

type Client struct {
	client    *httpclient.Client
	batchSize int
}

func New(cfg *config.FeedProviderConfig, logger logger.Logger) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil cfg")
	}

	client, err := httpclient.NewClient(cfg.Address, cfg.Timeout, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize http client: %w", err)
	}

	client.AddCircuitBreaker(providerName, &cfg.CircuitBreaker)

	batchSize := defaultBatchSize
	if cfg.BatchSize > 0 {
		batchSize = cfg.BatchSize
	}

	return &Client{
		client:    client,
		batchSize: batchSize,
	}, nil
}

func (h *Client) GetLatestNews(ctx context.Context) (*NewListInformation, error) {
	params := url.Values{
		countParam: {strconv.Itoa(h.batchSize)},
	}
	path := fmt.Sprintf("%s/%s?%s", basicAPIPath, getLatestNewsPath, params.Encode())
	req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	var resp NewListInformation
	if err = h.client.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to make a request: %w", err)
	}

	return &resp, nil
}
