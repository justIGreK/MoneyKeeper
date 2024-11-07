package mongorep

import (
	"budget/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Client) *UserRepo {
	return &UserRepo{
		collection: db.Database(dbname).Collection(userCollection),
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) (string, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *UserRepo) GetUser(ctx context.Context, userID string) (*models.User, error) {
	oid, err := convertToObjectIDs(userID)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": oid[0]}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, err
}
