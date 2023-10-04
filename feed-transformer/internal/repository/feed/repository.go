package feed

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vladislav-chunikhin/lib-go/pkg/logger"
	"github.com/vladislav-chunikhin/lib-go/pkg/mongodb"
)

const (
	dbName         = "sport_news"
	collectionName = "articles"
)

type Article struct {
	ID          string   `bson:"_id" json:"id"`
	ExternalID  int      `bson:"externalId" json:"externalId"`
	TeamID      *string  `bson:"teamId" json:"teamId"`
	OptaMatchID *string  `bson:"optaMatchId" json:"optaMatchId"`
	Title       *string  `bson:"title" json:"title"`
	Type        []string `bson:"type" json:"type"`
	Teaser      *string  `bson:"teaser" json:"teaser"`
	Content     *string  `bson:"content" json:"content"`
	URL         *string  `bson:"url" json:"url"`
	ImageURL    *string  `bson:"imageUrl" json:"imageUrl"`
	GalleryURLs *string  `bson:"galleryUrls" json:"galleryUrls"`
	VideoURL    *string  `bson:"videoUrl" json:"videoUrl"`
	Updated     *string  `bson:"updated" json:"updated"`
	Published   *string  `bson:"published" json:"published"`
}

type Repository struct {
	mongo  *mongodb.Client
	logger logger.Logger
}

func NewRepository(ctx context.Context, mongo *mongodb.Client, logger logger.Logger) (*Repository, error) {
	if err := createIndexes(ctx, mongo); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return &Repository{mongo: mongo, logger: logger}, nil
}

func (r *Repository) GetByExternalIDs(ctx context.Context, IDs []int, fields ...string) ([]*Article, error) {
	collection := r.mongo.GetCollection(dbName, collectionName)

	filter := bson.M{"externalId": bson.M{"$in": IDs}}

	projection := bson.M{}
	for _, field := range fields {
		projection[field] = 1
	}

	cur, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		if err = cur.Close(ctx); err != nil {
			r.logger.Warnf("failed to close mongo cursor")
		}
	}(cur, ctx)

	var articles []*Article
	for cur.Next(ctx) {
		var article Article
		if err = cur.Decode(&article); err != nil {
			return nil, err
		}
		articles = append(articles, &article)
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *Repository) UpdateOrCreateArticles(ctx context.Context, articles []*Article) error {
	collection := r.mongo.GetCollection(dbName, collectionName)

	var writeModels []mongo.WriteModel

	for _, article := range articles {
		filter := bson.M{filterKey: article.ExternalID}
		update := bson.M{
			"$set": bson.M{
				"teamId":      article.TeamID,
				"optaMatchId": article.OptaMatchID,
				"title":       article.Title,
				"type":        article.Type,
				"teaser":      article.Teaser,
				"content":     article.Content,
				"url":         article.URL,
				"imageUrl":    article.ImageURL,
				"galleryUrls": article.GalleryURLs,
				"videoUrl":    article.VideoURL,
				"updated":     article.Updated,
				"published":   article.Published,
			},
		}
		writeModels = append(writeModels, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true))
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := collection.BulkWrite(ctx, writeModels, opts)
	if err != nil {
		return err
	}

	return nil
}
