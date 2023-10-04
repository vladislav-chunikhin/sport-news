package mapper

import (
	"strings"
	"time"

	"github.com/vladislav-chunikhin/lib-go/pkg/logger"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/provider/htafc"
)

const (
	isPublished      = "true"
	inputDateLayout  = "2006-01-02 15:04:05"
	outputDateLayout = "2006-01-02T15:04:05.999Z"
)

type ArticleItems struct {
	Articles []*ArticleItem `json:"articles"`
}

type ArticleItem struct {
	ID        int    `json:"id"`
	Published string `json:"published"`
	Updated   string `json:"updated"`
}

func ToArticleItems(info *htafc.NewListInformation, logger logger.Logger) *ArticleItems {
	if info == nil {
		return nil
	}

	articlesItems := make([]*ArticleItem, 0, len(info.NewsletterNews.NewsletterNews))
	for _, article := range info.NewsletterNews.NewsletterNews {
		articleItem, err := toArticleItem(&article)
		if err != nil {
			logger.Errorf("failed to map an article in xml format to the article item: %v", err)
			continue
		}

		if articleItem != nil {
			articlesItems = append(articlesItems, articleItem)
		}
	}

	return &ArticleItems{Articles: articlesItems}
}

func toArticleItem(articleItem *htafc.NewsletterNewsItem) (*ArticleItem, error) {
	if articleItem == nil {
		return nil, nil
	}
	if strings.ToLower(articleItem.IsPublished) != isPublished {
		return nil, nil
	}

	updateAsTime, err := time.Parse(inputDateLayout, articleItem.LastUpdateDate)
	if err != nil {
		return nil, err
	}
	updateAsStr := updateAsTime.Format(outputDateLayout)

	var publishedAsTime time.Time
	publishedAsTime, err = time.Parse(inputDateLayout, articleItem.PublishDate)
	if err != nil {
		return nil, err
	}
	publishedAsStr := publishedAsTime.Format(outputDateLayout)

	return &ArticleItem{
		ID:        articleItem.NewsArticleID,
		Updated:   updateAsStr,
		Published: publishedAsStr,
	}, nil
}
