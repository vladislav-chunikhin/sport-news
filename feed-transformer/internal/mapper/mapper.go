package mapper

import (
	"strings"
	"time"

	"github.com/vladislav-chunikhin/feed-transformer/internal/provider/htafc"
	feedRepoPkg "github.com/vladislav-chunikhin/feed-transformer/internal/repository/feed"
)

const (
	isPublished      = "true"
	inputDateLayout  = "2006-01-02 15:04:05"
	outputDateLayout = "2006-01-02T15:04:05.999Z"
	htafcType        = "htafc"
)

func ToArticle(articleItem *htafc.NewsArticleInformation, msgType string) (*feedRepoPkg.Article, error) {
	if articleItem == nil {
		return nil, nil
	}
	if strings.ToLower(articleItem.NewsArticle.IsPublished) != isPublished {
		return nil, nil
	}

	var result feedRepoPkg.Article

	updateAsTime, err := time.Parse(inputDateLayout, articleItem.NewsArticle.LastUpdateDate)
	if err != nil {
		return nil, err
	}
	updateAsStr := updateAsTime.Format(outputDateLayout)

	var publishedAsTime time.Time
	publishedAsTime, err = time.Parse(inputDateLayout, articleItem.NewsArticle.PublishDate)
	if err != nil {
		return nil, err
	}
	publishedAsStr := publishedAsTime.Format(outputDateLayout)

	teamID := ""
	if msgType == htafcType {
		teamID = "t94"
	}

	result.ExternalID = articleItem.NewsArticle.NewsArticleID
	result.TeamID = &teamID
	result.OptaMatchID = &articleItem.NewsArticle.OptaMatchId
	result.Title = &articleItem.NewsArticle.Title
	result.Type = []string{articleItem.NewsArticle.Taxonomies}
	result.Teaser = &articleItem.NewsArticle.TeaserText
	result.Content = &articleItem.NewsArticle.BodyText
	result.URL = &articleItem.NewsArticle.ArticleURL
	result.ImageURL = &articleItem.NewsArticle.ThumbnailImageURL
	result.GalleryURLs = &articleItem.NewsArticle.GalleryImageURLs
	result.VideoURL = &articleItem.NewsArticle.VideoURL
	result.Updated = &updateAsStr
	result.Published = &publishedAsStr

	return &result, nil
}
