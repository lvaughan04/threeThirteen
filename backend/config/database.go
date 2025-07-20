package config

import (
	"log"
	"os"

	"context"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *mongo.Client
}

func NewDatabase() (*Database, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, log.Output(2, "MONGO_URI not set in .env file")
	}

	clientOptions := options.Client().ApplyURI(mongoURI).SetConnectTimeout(10 * time.Second)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB successfully")
	return &Database{Client: client}, nil
}

func (db *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return db.Client.Disconnect(ctx)
}
