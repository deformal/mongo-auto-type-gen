package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SampleDocuments(ctx context.Context, coll *mongo.Collection, n int) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: n}}}},
	}

	cur, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return FindFirstN(ctx, coll, n)
	}
	defer cur.Close(ctx)

	var docs []bson.M
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func FindFirstN(ctx context.Context, coll *mongo.Collection, n int) ([]bson.M, error) {

	findOpts := &options.FindOptions{
		Limit: ptrInt64(int64(n)),
		Sort:  bson.D{{Key: "_id", Value: -1}},
	}

	cur, err := coll.Find(ctx, bson.M{}, findOpts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var docs []bson.M
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func ptrInt64(v int64) *int64 { return &v }
