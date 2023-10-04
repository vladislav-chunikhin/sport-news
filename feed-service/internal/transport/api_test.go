package transport

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/vladislav-chunikhin/lib-go/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"

	feedRepoPkg "github.com/vladislav-chunikhin/feed-service/internal/repository/feed"
	feedServicePkg "github.com/vladislav-chunikhin/feed-service/internal/service/feed"
	"github.com/vladislav-chunikhin/feed-service/internal/transport/mocks"
)

func TestAPI_GetByID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedService := mocks.NewMockFeedService(ctrl)
	expectedArticle := &feedRepoPkg.Article{}

	mockFeedService.EXPECT().
		GetByID(gomock.Any(), gomock.Any()).
		Return(expectedArticle, nil)

	log, err := logger.New(logger.DebugLevel)
	require.NoError(t, err)

	api := NewAPI(mockFeedService, log)

	r := chi.NewRouter()
	r.Route(pattern, func(r chi.Router) {
		r.Get("/{id}", api.GetByID)
	})

	path := "/feed/v1/news/64dcc2326c1e91ddcad83046"
	req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestAPI_GetLatestArticles_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedService := mocks.NewMockFeedService(ctrl)
	expectedResult := &feedServicePkg.ArticlePaginationResult{}

	mockFeedService.EXPECT().
		GetLatestArticles(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(expectedResult, nil)

	log, err := logger.New(logger.DebugLevel)
	require.NoError(t, err)

	api := NewAPI(mockFeedService, log)

	r := chi.NewRouter()
	r.Route(pattern, func(r chi.Router) {
		r.Get("/", api.GetLatestArticles)
	})

	path := "/feed/v1/news"
	req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestAPI_GetByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.MockLogger{}

	testCases := []struct {
		name               string
		feedService        func() FeedService
		path               string
		expectedHTTPStatus int
	}{
		{
			name: "invalid ID",
			feedService: func() FeedService {
				return nil
			},
			path:               "/feed/v1/news/64dcc",
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "service err",
			feedService: func() FeedService {
				mockFeedService := mocks.NewMockFeedService(ctrl)

				mockFeedService.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return mockFeedService
			},
			path:               "/feed/v1/news/64dcc2326c1e91ddcad83046",
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name: "not found by id",
			feedService: func() FeedService {
				mockFeedService := mocks.NewMockFeedService(ctrl)

				mockFeedService.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, mongo.ErrNoDocuments)

				return mockFeedService
			},
			path:               "/feed/v1/news/64dcc2326c1e91ddcad83046",
			expectedHTTPStatus: http.StatusNotFound,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(tt.feedService(), &log)

			r := chi.NewRouter()
			r.Route(pattern, func(r chi.Router) {
				r.Get("/{id}", api.GetByID)
			})

			var req *http.Request
			var err error
			req, err = http.NewRequest(http.MethodGet, tt.path, http.NoBody)
			require.NoError(t, err)
			require.NotNil(t, req)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedHTTPStatus, rr.Code)
		})
	}
}

func TestAPI_GetLatestArticles_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.MockLogger{}

	testCases := []struct {
		name               string
		feedService        func() FeedService
		path               string
		expectedHTTPStatus int
	}{
		{
			name: "invalid cursor",
			feedService: func() FeedService {
				return nil
			},
			path:               "/feed/v1/news?cursor=123",
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "invalid limit",
			feedService: func() FeedService {
				return nil
			},
			path:               "/feed/v1/news?limit=abc",
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "service err",
			feedService: func() FeedService {
				mockFeedService := mocks.NewMockFeedService(ctrl)

				mockFeedService.EXPECT().
					GetLatestArticles(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return mockFeedService
			},
			path:               "/feed/v1/news",
			expectedHTTPStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(tt.feedService(), &log)

			r := chi.NewRouter()
			r.Route(pattern, func(r chi.Router) {
				r.Get("/", api.GetLatestArticles)
			})

			var req *http.Request
			var err error
			req, err = http.NewRequest(http.MethodGet, tt.path, http.NoBody)
			require.NoError(t, err)
			require.NotNil(t, req)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedHTTPStatus, rr.Code)
		})
	}
}
