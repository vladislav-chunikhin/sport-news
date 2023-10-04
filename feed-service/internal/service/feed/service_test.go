package feed

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	feedRepoPkg "github.com/vladislav-chunikhin/feed-service/internal/repository/feed"
	"github.com/vladislav-chunikhin/feed-service/internal/service/feed/mocks"
)

func TestService_GetByID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	expectedArticle := &feedRepoPkg.Article{}

	mockRepo.EXPECT().
		GetByID(gomock.Any(), gomock.Any()).
		Return(expectedArticle, nil)

	service := NewService(mockRepo)

	article, err := service.GetByID(context.Background(), primitive.NewObjectID())
	require.NoError(t, err)
	require.Equal(t, expectedArticle, article)
}

func TestService_GetByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		GetByID(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("some error"))

	service := NewService(mockRepo)

	article, err := service.GetByID(context.Background(), primitive.NewObjectID())
	require.Error(t, err)
	require.Nil(t, article)
}

func TestService_GetLatestArticles_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	expectedArticle := &feedRepoPkg.Article{}

	mockRepo.EXPECT().
		GetLatestArticles(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*feedRepoPkg.Article{expectedArticle}, nil)

	service := NewService(mockRepo)

	result, err := service.GetLatestArticles(context.Background(), "", 1)
	require.NoError(t, err)
	require.Equal(t, 1, result.TotalItems)
	require.Equal(t, defaultSortType, result.Sort)
}
