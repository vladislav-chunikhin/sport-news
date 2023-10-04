package transport

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	cursorKey = "cursor"
	limitKey  = "limit"
	IDKey     = "id"

	defaultLimit  = 50
	maxLimitValue = 200
)

var (
	ErrInvalidCursor  = fmt.Errorf("invalid cursor value")
	ErrInvalidLimit   = fmt.Errorf("invalid limit value")
	ErrNilResult      = fmt.Errorf("nil result")
	ErrIDNotSpecified = fmt.Errorf("ID must be specified")
	ErrInvalidID      = fmt.Errorf("invalid ID value")
)

func (a *API) GetLatestArticles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cursor := r.URL.Query().Get(cursorKey)
	if cursor != "" {
		if _, err := time.Parse(dateFormat, cursor); err != nil {
			a.logger.Errorf("%s, invalid cursor: %v", r.RequestURI, err)
			errorJSON(w, r, http.StatusBadRequest, ErrInvalidCursor)
			return
		}
	}

	limitStr := r.URL.Query().Get(limitKey)
	limit := defaultLimit
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			a.logger.Errorf("%s, invalid limit parameter: %v", r.RequestURI, err)
			errorJSON(w, r, http.StatusBadRequest, ErrInvalidLimit)
			return
		}
	}

	if limit > maxLimitValue {
		a.logger.Warnf("limit value too high, default value will be used")
		limit = defaultLimit
	}

	result, err := a.feedService.GetLatestArticles(ctx, cursor, limit)
	if err != nil {
		a.logger.Errorf("%s, failed to get articles: %v", r.RequestURI, err)
		errorJSON(w, r, http.StatusInternalServerError, err)
		return
	}

	if result == nil {
		a.logger.Errorf("%s, nil result", r.RequestURI)
		errorJSON(w, r, http.StatusInternalServerError, ErrNilResult)
		return
	}

	okPaginationJSON(w, r, result.Articles, result.Sort, result.TotalItems)
}

func (a *API) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ID := chi.URLParam(r, IDKey)
	if ID == "" {
		a.logger.Errorf("ID must be specified")
		errorJSON(w, r, http.StatusBadRequest, ErrIDNotSpecified)
		return
	}

	var objID primitive.ObjectID
	var err error
	if objID, err = primitive.ObjectIDFromHex(ID); err != nil {
		a.logger.Errorf("%s, invalid ID value: %v", r.RequestURI, err)
		errorJSON(w, r, http.StatusBadRequest, ErrInvalidID)
		return
	}

	result, err := a.feedService.GetByID(ctx, objID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			errorJSON(w, r, http.StatusNotFound, err)
			return
		}
		a.logger.Errorf("%s, failed to get article by ID: %v", r.RequestURI, err)
		errorJSON(w, r, http.StatusInternalServerError, err)
		return
	}

	if result == nil {
		a.logger.Errorf("%s, nil result", r.RequestURI)
		errorJSON(w, r, http.StatusInternalServerError, ErrNilResult)
		return
	}

	okJSON(w, r, result)
}
