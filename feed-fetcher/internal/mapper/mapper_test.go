package mapper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vladislav-chunikhin/feed-fetcher/internal/mapper/mocks"
	"github.com/vladislav-chunikhin/feed-fetcher/internal/provider/htafc"
)

func TestToArticleItems(t *testing.T) {
	mockLogger := mocks.NewMockLogger()

	info := &htafc.NewListInformation{
		ClubName:       "Sample Club",
		ClubWebsiteURL: "https://sampleclub.com",
		NewsletterNews: htafc.NewsletterNewsList{
			NewsletterNews: []htafc.NewsletterNewsItem{
				{
					NewsArticleID:  1,
					PublishDate:    "2023-08-01 12:00:00",
					LastUpdateDate: "2023-08-02 10:00:00",
					IsPublished:    "True",
				},
				{
					NewsArticleID:  2,
					PublishDate:    "2023-08-02 14:00:00",
					LastUpdateDate: "2023-08-03 11:30:00",
					IsPublished:    "True",
				},
			},
		},
	}

	expectedArticleItems := &ArticleItems{
		Articles: []*ArticleItem{
			{
				ID:        1,
				Published: "2023-08-01T12:00:00Z",
				Updated:   "2023-08-02T10:00:00Z",
			},
			{
				ID:        2,
				Published: "2023-08-02T14:00:00Z",
				Updated:   "2023-08-03T11:30:00Z",
			},
		},
	}

	result := ToArticleItems(info, mockLogger)

	require.Equal(t, len(expectedArticleItems.Articles), len(result.Articles))

	for i, expectedArticle := range expectedArticleItems.Articles {
		require.Equal(t, expectedArticle.ID, result.Articles[i].ID)
		require.Equal(t, expectedArticle.Published, result.Articles[i].Published)
		require.Equal(t, expectedArticle.Updated, result.Articles[i].Updated)
	}
}

func TestToArticleItems_Nil_Input(t *testing.T) {
	require.Nil(t, ToArticleItems(nil, mocks.NewMockLogger()))
}

func TestToArticleItem(t *testing.T) {
	testCases := []struct {
		name           string
		input          *htafc.NewsletterNewsItem
		expectedResult *ArticleItem
		expectedError  error
	}{
		{
			name:           "nil input",
			input:          nil,
			expectedResult: nil,
			expectedError:  nil,
		},
		{
			name: "article item is not published",
			input: &htafc.NewsletterNewsItem{
				NewsArticleID:  1,
				PublishDate:    "2023-08-01 12:00:00",
				LastUpdateDate: "2023-08-02 10:00:00",
				IsPublished:    "False",
			},
			expectedResult: nil,
			expectedError:  nil,
		},
		{
			name: "invalid last update date",
			input: &htafc.NewsletterNewsItem{
				NewsArticleID:  1,
				PublishDate:    "2023-08-01 12:00:00",
				LastUpdateDate: "2023-08",
				IsPublished:    "True",
			},
			expectedResult: nil,
			expectedError:  fmt.Errorf("parsing time \"2023-08\" as \"2006-01-02 15:04:05\": cannot parse \"\" as \"-\""),
		},
		{
			name: "invalid publish date",
			input: &htafc.NewsletterNewsItem{
				NewsArticleID:  1,
				PublishDate:    "2023-08",
				LastUpdateDate: "2023-08-02 10:00:00",
				IsPublished:    "True",
			},
			expectedResult: nil,
			expectedError:  fmt.Errorf("parsing time \"2023-08\" as \"2006-01-02 15:04:05\": cannot parse \"\" as \"-\""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := toArticleItem(tc.input)
			require.Equal(t, tc.expectedResult, result)
			if err != nil {
				require.NotNil(t, tc.expectedError)
				require.Equal(t, tc.expectedError.Error(), err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
