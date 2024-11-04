package mongorep

import (
	"context"
	"errors"
	"log"

	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbname         = "budgetdb"
	userCollection = "users"
	budgetCollection = "budgets"
	transactionCollection = "transactions"
)

func CreateMongoClient(ctx context.Context) *mongo.Client {
	dbURI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB is not connected: %v", err)
	}
	return client
}

func convertToObjectIDs(ids ...string) ([]primitive.ObjectID, error) {
	objectIDs := make([]primitive.ObjectID, 0, len(ids))

	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New("InvalidID: " + id)
		}
		objectIDs = append(objectIDs, oid)
	}

	return objectIDs, nil
}
