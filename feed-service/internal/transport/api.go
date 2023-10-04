package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/vladislav-chunikhin/lib-go/pkg/logger"

	"github.com/vladislav-chunikhin/feed-service/internal/model"
	feedRepoPkg "github.com/vladislav-chunikhin/feed-service/internal/repository/feed"
	feedServicePkg "github.com/vladislav-chunikhin/feed-service/internal/service/feed"
)

const (
	pattern   = "/feed/v1/news"
	errStatus = "error"
	okStatus  = "success"

	dateFormat = "2006-01-02T15:04:05.999Z"
)

//go:generate mockgen -source=api.go -destination=mocks/feed_service.go -package=mocks
type FeedService interface {
	GetByID(ctx context.Context, ID primitive.ObjectID) (*feedRepoPkg.Article, error)
	GetLatestArticles(ctx context.Context, cursor string, limit int) (*feedServicePkg.ArticlePaginationResult, error)
}

type API struct {
	feedService FeedService
	logger      logger.Logger
}

func NewAPI(feedService FeedService, logger logger.Logger) *API {
	return &API{feedService: feedService, logger: logger}
}

func (a *API) Router() chi.Router {
	r := chi.NewRouter()

	r.Route(pattern, func(r chi.Router) {
		r.Get("/", a.GetLatestArticles)
		r.Get("/{id}", a.GetByID)
	})

	return r
}

func errorJSON(w http.ResponseWriter, r *http.Request, code int, err error) {
	var resp model.ErrorResult

	resp.Message = err.Error()
	resp.Status = errStatus
	resp.Metadata.CreatedAt = time.Now().Format(dateFormat)

	render.Status(r, code)
	render.JSON(w, r, &resp)
}

func okJSON(w http.ResponseWriter, r *http.Request, obj any) {
	if obj == nil {
		obj = struct{}{}
	}

	var resp model.OkResult

	resp.Data = obj
	resp.Status = okStatus
	resp.Metadata.CreatedAt = time.Now().Format(dateFormat)

	render.JSON(w, r, &resp)
}

func okPaginationJSON(w http.ResponseWriter, r *http.Request, obj any, sort string, totalItems int) {
	if obj == nil {
		obj = struct{}{}
	}

	var resp model.OkPaginationResult

	resp.Data = obj
	resp.Status = okStatus
	resp.Metadata.CreatedAt = time.Now().Format(dateFormat)
	resp.Metadata.Sort = sort
	resp.Metadata.TotalItems = totalItems

	render.JSON(w, r, &resp)
}
