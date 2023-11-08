package fetcher

import (
	"context"
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/fetcher/mocks"
	"github.com/vladislav-chunikhin/feed-fetcher/internal/provider/htafc"
)

func TestFetcher_Fetch(t *testing.T) {
	ctx := context.TODO()
	mockLogger := mocks.NewMockLogger()

	testCases := []struct {
		name           string
		producer       Producer
		provider       FeedProvider
		expectedErrMsg string
	}{
		{
			name: "success",
			producer: func() Producer {
				producer := mocks.NewProducer(t)
				producer.EXPECT().
					PublishHtafcFeed(ctx, mock.Anything).
					Return(nil)
				return producer
			}(),
			provider: func() FeedProvider {
				provider := mocks.NewFeedProvider(t)
				provider.EXPECT().
					GetLatestNews(ctx).
					Return(getNewListInformation(), nil)
				return provider
			}(),
		},
		{
			name:     "provider err",
			producer: mocks.NewProducer(t),
			provider: func() FeedProvider {
				provider := mocks.NewFeedProvider(t)
				provider.EXPECT().
					GetLatestNews(ctx).
					Return(nil, fmt.Errorf("provider conn error"))
				return provider
			}(),
			expectedErrMsg: "failed to get latest news: provider conn error",
		},
		{
			name:     "nil last news",
			producer: mocks.NewProducer(t),
			provider: func() FeedProvider {
				provider := mocks.NewFeedProvider(t)
				provider.EXPECT().
					GetLatestNews(ctx).
					Return(nil, nil)
				return provider
			}(),
			expectedErrMsg: "nil result from feed client",
		},
		{
			name:     "empty last news",
			producer: mocks.NewProducer(t),
			provider: func() FeedProvider {
				provider := mocks.NewFeedProvider(t)
				expectedDataProvider := getNewListInformation()
				expectedDataProvider.NewsletterNews.NewsletterNews = []htafc.NewsletterNewsItem{}

				provider.EXPECT().
					GetLatestNews(ctx).
					Return(expectedDataProvider, nil)
				return provider
			}(),
			expectedErrMsg: "no feeds from provider",
		},
		{
			name: "publish error",
			producer: func() Producer {
				producer := mocks.NewProducer(t)
				producer.EXPECT().
					PublishHtafcFeed(ctx, mock.Anything).
					Return(fmt.Errorf("some error"))
				return producer
			}(),
			provider: func() FeedProvider {
				provider := mocks.NewFeedProvider(t)
				provider.EXPECT().
					GetLatestNews(ctx).
					Return(getNewListInformation(), nil)
				return provider
			}(),
			expectedErrMsg: "failed to publish message: some error",
		},
		{
			name:           "nil producer",
			producer:       nil,
			provider:       mocks.NewFeedProvider(t),
			expectedErrMsg: "nil producer",
		},
		{
			name:           "nil provider",
			producer:       mocks.NewProducer(t),
			provider:       nil,
			expectedErrMsg: "nil htafc feed provider",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fetcher, err := NewFetcher(tc.producer, tc.provider, mockLogger)
			if tc.expectedErrMsg != "" && err != nil {
				require.Error(t, err)
				require.Nil(t, fetcher)
				require.Equal(t, tc.expectedErrMsg, err.Error())
				return
			}

			require.NoError(t, err)
			err = fetcher.Fetch(ctx)
			if tc.expectedErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.expectedErrMsg, err.Error())
			}
		})
	}
}

func getNewListInformation() *htafc.NewListInformation {
	return &htafc.NewListInformation{
		XMLName: xml.Name{
			Local: "NewListInformation",
		},
		ClubName:       "Sample Club",
		ClubWebsiteURL: "https://www.sampleclub.com",
		NewsletterNews: htafc.NewsletterNewsList{
			XMLName: xml.Name{
				Local: "NewsletterNewsItems",
			},
			NewsletterNews: []htafc.NewsletterNewsItem{
				{
					XMLName: xml.Name{
						Local: "NewsletterNewsItem",
					},
					ArticleURL:        "https://www.sampleclub.com/news/article1",
					NewsArticleID:     1,
					PublishDate:       "2006-01-02 15:04:05",
					Taxonomies:        "Sports",
					TeaserText:        "Sample teaser text for article 1",
					ThumbnailImageURL: "https://www.sampleclub.com/images/thumbnail1.jpg",
					Title:             "Sample Article 1",
					OptaMatchId:       "12345",
					LastUpdateDate:    "2006-01-02 15:04:05",
					IsPublished:       "true",
				},
				{
					XMLName: xml.Name{
						Local: "NewsletterNewsItem",
					},
					ArticleURL:        "https://www.sampleclub.com/news/article2",
					NewsArticleID:     2,
					PublishDate:       "2023-11-06",
					Taxonomies:        "Entertainment",
					TeaserText:        "Sample teaser text for article 2",
					ThumbnailImageURL: "https://www.sampleclub.com/images/thumbnail2.jpg",
					Title:             "Sample Article 2",
					OptaMatchId:       "67890",
					LastUpdateDate:    "2023-11-06",
					IsPublished:       "true",
				},
			},
		},
	}
}
