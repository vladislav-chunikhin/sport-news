package feed

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vladislav-chunikhin/lib-go/pkg/mongodb"
)

const (
	nameKey            = "name"
	publishedIndexName = "published_index"
)

func createIndexes(ctx context.Context, mongoClient *mongodb.Client) error {
	collection := mongoClient.GetCollection(dbName, collectionName)
	indexes, err := collection.Indexes().List(ctx)
	if err != nil {
		return err
	}
	var existingIndexes []bson.M
	if err = indexes.All(ctx, &existingIndexes); err != nil {
		return err
	}

	hasPublishedIndex := false
	for _, index := range existingIndexes {
		if indexName, ok := index[nameKey]; ok {
			if indexName == publishedIndexName {
				hasPublishedIndex = true
			}
		}
	}

	if !hasPublishedIndex {
		index := mongo.IndexModel{
			Keys:    bson.M{defaultSortKey: -1},
			Options: options.Index().SetName(publishedIndexName),
		}
		_, err = collection.Indexes().CreateOne(ctx, index)
		if err != nil {
			return err
		}
	}
	return nil
}
