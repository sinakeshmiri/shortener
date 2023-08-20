package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Adapter implements the DbPort interface
type Adapter struct {
	client *mongo.Client
}

// NewAdapter creates a new Adapter
func NewAdapter(mongouri string) (*Adapter, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(mongouri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return &Adapter{client: client}, nil
}

func (da Adapter) CloseDbConnection() {
	if err := da.client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func (da Adapter) AddURL(url, urlID, username string) error {
	collection := da.client.Database("urlshortener").Collection("urls")

	// Convert "hits" to an integer value (0)
	hits := 0

	_, err := collection.InsertOne(context.Background(), map[string]interface{}{
		"id":       urlID,
		"url":      url,
		"username": username,
		"hits":     hits, // Insert as an integer
	})

	if err != nil {

		return err
	}

	return nil
}

func (da Adapter) DeleteURL(id, username string) error {
	collection := da.client.Database("urlshortener").Collection("urls")
	_, err := collection.DeleteOne(context.Background(), map[string]string{"id": id, "username": username})
	if err != nil {
		return err
	}
	return nil
}
func (da Adapter) GetHits(username string) (map[string]int, error) {
	collection := da.client.Database("urlshortener").Collection("urls")
	cursor, err := collection.Find(context.Background(), map[string]string{"username": username})
	if err != nil {
		return nil, err
	}
	var url map[string]any
	metrics := make(map[string]int)
	for cursor.Next(context.Background()) {
		if err = cursor.Decode(&url); err != nil {
			return nil, err
		}
		metrics[url["id"].(string)] = 0
	}
	return metrics, nil
}

func (da Adapter) GetURL(id string) (string, error) {
	collection := da.client.Database("urlshortener").Collection("urls")
	var url map[string]any
	err := collection.FindOne(context.Background(), map[string]string{"id": id}).Decode(&url)
	if err != nil {
		return "", err
	}
	return url["url"].(string), nil
}
func (da Adapter) AddHit(id string) error {
	collection := da.client.Database("urlshortener").Collection("urls")

	filter := bson.M{"id": id}
	update := bson.M{"$inc": bson.M{"hits": 1}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}
