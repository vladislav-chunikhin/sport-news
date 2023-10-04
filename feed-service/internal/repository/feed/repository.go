package feed

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vladislav-chunikhin/lib-go/pkg/logger"
	"github.com/vladislav-chunikhin/lib-go/pkg/mongodb"
)

const (
	dbName         = "sport_news"
	collectionName = "articles"

	defaultSortKey = "published"
	idKey          = "_id"
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
	mongoClient *mongodb.Client
	logger      logger.Logger
}

func NewRepository(ctx context.Context, mongoClient *mongodb.Client, logger logger.Logger) (*Repository, error) {
	if err := createIndexes(ctx, mongoClient); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return &Repository{mongoClient: mongoClient, logger: logger}, nil
}

func (r *Repository) GetByID(ctx context.Context, ID primitive.ObjectID) (*Article, error) {
	filter := bson.M{idKey: ID}
	var article Article

	if err := r.mongoClient.GetCollection(dbName, collectionName).FindOne(ctx, filter).Decode(&article); err != nil {
		return nil, err
	}

	return &article, nil
}

func (r *Repository) GetLatestArticles(ctx context.Context, cursor string, limit int) ([]*Article, error) {
	collection := r.mongoClient.GetCollection(dbName, collectionName)

	opts := options.Find()
	opts.SetSort(bson.D{{defaultSortKey, -1}})
	opts.SetLimit(int64(limit))

	filter := bson.M{}
	if cursor != "" {
		filter = bson.M{defaultSortKey: bson.M{"$lt": cursor}}
	}

	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error fetching articles: %w", err)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		if err = cur.Close(ctx); err != nil {
			r.logger.Warnf("failed to close mongo cursor")
		}
	}(cur, ctx)

	articles := make([]*Article, 0, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		var article Article
		if err = cur.Decode(&article); err != nil {
			return nil, fmt.Errorf("error decoding article: %w", err)
		}
		articles = append(articles, &article)
	}

	if err = cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return articles, nil
}
