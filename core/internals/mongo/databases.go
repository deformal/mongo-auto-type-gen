package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func ListDatabases(ctx context.Context, client *mongo.Client) ([]string, error) {
	list, err := client.ListDatabaseNames(ctx, map[string]any{})
	if err != nil {
		fmt.Println("Mongo connection error while listing db's")
		fmt.Println(err)
		return nil, err
	}
	var databaseNames = []string{}
	for _, db := range list {
		if db == "admin" || db == "local" || db == "config" {
			continue
		}
		databaseNames = append(databaseNames, db)
	}
	return databaseNames, nil
}
