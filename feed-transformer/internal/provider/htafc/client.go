package htafc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/vladislav-chunikhin/lib-go/pkg/httpclient"
	"github.com/vladislav-chunikhin/lib-go/pkg/logger"

	"github.com/vladislav-chunikhin/feed-transformer/internal/config"
)

const (
	basicAPIPath           = "/api/incrowd"
	getNewsContentByIDPath = "getnewsarticleinformation"

	providerName = "htafc"

	idParam = "id"
)

type Htafc struct {
	client *httpclient.Client
}

func New(cfg *config.FeedProviderConfig, limit int, logger logger.Logger) (*Htafc, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil cfg")
	}

	client, err := httpclient.NewClient(cfg.Address, cfg.Timeout, logger)
	if err != nil {
		return nil, err
	}

	client.AddCircuitBreaker(providerName, &cfg.CircuitBreaker)
	client.AddLimiter(limit)

	return &Htafc{
		client: client,
	}, nil
}

func (h *Htafc) GetNewsContentByID(ctx context.Context, ID int) (*NewsArticleInformation, error) {
	params := url.Values{
		idParam: {strconv.Itoa(ID)},
	}
	path := fmt.Sprintf("%s/%s?%s", basicAPIPath, getNewsContentByIDPath, params.Encode())
	req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	var resp NewsArticleInformation
	if err = h.client.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to make a request: %w", err)
	}

	return &resp, nil
}
