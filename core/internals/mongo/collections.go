package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func ListCollections(ctx context.Context, db *mongo.Database) ([]string, error) {
	return db.ListCollectionNames(ctx, map[string]any{})
}
