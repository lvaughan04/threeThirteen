package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	var user models.User

	filter := bson.M{"_id": userID}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	filter := bson.M{"username": username}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetAllUsers retrieves all users from the database
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(ctx context.Context, userID primitive.ObjectID, user *models.User) error {
	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": user}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) UserExists(ctx context.Context, username, email string) (bool, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": email},
		},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return count > 0, nil
}
