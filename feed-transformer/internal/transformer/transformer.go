package transformer

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/vladislav-chunikhin/feed-transformer/internal/provider/htafc"
	"github.com/vladislav-chunikhin/lib-go/pkg/logger"

	feedRepoPkg "github.com/vladislav-chunikhin/feed-transformer/internal/repository/feed"
)

type Cache interface {
	Get(ctx context.Context, key string, ptrValue any) error
	Set(ctx context.Context, key string, ptrValue any, expires time.Duration) error
}

type Consumer interface {
	ConsumeHtafcFeed(ctx context.Context, handler func(amqp.Delivery) error) error
}

type FeedClient interface {
	GetNewsContentByID(ctx context.Context, ID int) (*htafc.NewsArticleInformation, error)
}

type FeedRepository interface {
	GetByExternalIDs(ctx context.Context, IDs []int, fields ...string) ([]*feedRepoPkg.Article, error)
	UpdateOrCreateArticles(ctx context.Context, articles []*feedRepoPkg.Article) error
}

type ArticleItems struct {
	Articles   []*ArticleItem `json:"articles"`
	TotalItems int            `json:"totalItems"`
}

type ArticleItem struct {
	ID        int    `json:"id"`
	Published string `json:"published"`
	Updated   string `json:"updated"`
}

type Service struct {
	ctx            context.Context
	consumer       Consumer
	feedClient     FeedClient
	feedRepository FeedRepository
	cache          Cache
	logger         logger.Logger
}

func NewService(
	ctx context.Context,
	consumer Consumer,
	feedClient FeedClient,
	feedRepository FeedRepository,
	cache Cache,
	logger logger.Logger,
) *Service {
	return &Service{
		ctx:            ctx,
		consumer:       consumer,
		feedClient:     feedClient,
		feedRepository: feedRepository,
		cache:          cache,
		logger:         logger,
	}
}

func (s *Service) Run(ctx context.Context) (err error) {
	err = s.consumer.ConsumeHtafcFeed(ctx, s.handle)
	return
}
