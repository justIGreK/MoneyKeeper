package repository

import (
	"context"

	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
}

const (
	dbname         = "budgets"
	userCollection = "users"
)

func NewUserRepo(db *mongo.Client) *UserRepo {
	return &UserRepo{
		collection: db.Database(dbname).Collection(userCollection),
	}
}

func (r *UserRepo) AddUser(ctx context.Context, chatID, userID string) error {
	_, err := r.collection.InsertOne(ctx, bson.M{"chat_id": chatID, "user_id": userID})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetUserID(ctx context.Context, chatID string) (string, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}
	return user.UserID, nil
}
