package feed

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	feedRepoPkg "github.com/vladislav-chunikhin/feed-service/internal/repository/feed"
)

const defaultSortType = "-published"

type ArticlePaginationResult struct {
	Articles   []*feedRepoPkg.Article
	Sort       string
	TotalItems int
}

//go:generate mockgen -source=service.go -destination=mocks/feed_repository.go -package=mocks
type Repository interface {
	GetByID(ctx context.Context, ID primitive.ObjectID) (*feedRepoPkg.Article, error)
	GetLatestArticles(ctx context.Context, cursor string, limit int) ([]*feedRepoPkg.Article, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, ID primitive.ObjectID) (*feedRepoPkg.Article, error) {
	article, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get article by ID: %w", err)
	}
	return article, nil
}

func (s *Service) GetLatestArticles(ctx context.Context, cursor string, limit int) (*ArticlePaginationResult, error) {
	articles, err := s.repo.GetLatestArticles(ctx, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %w", err)
	}

	sort := defaultSortType
	if len(articles) == 0 {
		sort = ""
	}

	return &ArticlePaginationResult{
		Articles:   articles,
		Sort:       sort,
		TotalItems: len(articles),
	}, nil
}
