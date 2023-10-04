package mapper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vladislav-chunikhin/feed-transformer/internal/provider/htafc"
)

func TestToArticle(t *testing.T) {
	articleInfo := &htafc.NewsArticleInformation{
		ClubName: "Sample Club",
		NewsArticle: htafc.NewsArticle{
			IsPublished:       "True",
			NewsArticleID:     123,
			OptaMatchId:       "123456",
			Title:             "Sample Article",
			Taxonomies:        "Sports",
			TeaserText:        "Sample Teaser",
			BodyText:          "Sample Content",
			ArticleURL:        "http://example.com/article",
			ThumbnailImageURL: "http://example.com/thumbnail.jpg",
			GalleryImageURLs:  "http://example.com/image1.jpg,http://example.com/image2.jpg",
			VideoURL:          "http://example.com/video.mp4",
			LastUpdateDate:    "2023-08-17 12:00:00",
			PublishDate:       "2023-08-16 12:00:00",
		},
	}

	article, err := ToArticle(articleInfo, htafcType)

	require.NoError(t, err, "Expected no error")

	require.NotNil(t, article, "Expected a non-nil article")

	require.Equal(t, 123, article.ExternalID, "ExternalID should match")
	require.Equal(t, "t94", *article.TeamID, "TeamID should match")
	require.Equal(t, "123456", *article.OptaMatchID, "OptaMatchID should match")
	require.Equal(t, "Sample Article", *article.Title, "Title should match")
	require.Equal(t, []string{"Sports"}, article.Type, "Type should match")
	require.Equal(t, "Sample Teaser", *article.Teaser, "Teaser should match")
	require.Equal(t, "Sample Content", *article.Content, "Content should match")
	require.Equal(t, "http://example.com/article", *article.URL, "URL should match")
	require.Equal(t, "http://example.com/thumbnail.jpg", *article.ImageURL, "ImageURL should match")
	require.Equal(t, "http://example.com/image1.jpg,http://example.com/image2.jpg", *article.GalleryURLs, "GalleryURLs should match")
	require.Equal(t, "http://example.com/video.mp4", *article.VideoURL, "VideoURL should match")
}
