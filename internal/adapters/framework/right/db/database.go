package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

// Adapter implements the DbPort interface
type Adapter struct {
	client *mongo.Client
}

// NewAdapter creates a new Adapter
func NewAdapter(mongouri string) (*Adapter, error) {
	// Set client options
	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI("mongodb://mongo:27017")
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
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
	newUrl := urlStruct{
		ID:       urlID,
		URL:      url,
		Username: username,
		Hits:     0,
	}
	collection := da.client.Database("urlshortener").Collection("urls")
	_, err := collection.InsertOne(context.Background(), newUrl)
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
func (da Adapter) GetHits(username, id string) (map[string]int, error) {
	collection := da.client.Database("urlshortener").Collection("urls")
	if id != "" {
		var url urlStruct
		err := collection.FindOne(context.Background(), map[string]string{"id": id, "username": username}).Decode(&url)
		if err != nil {
			return nil, err
		}
		metric := make(map[string]int)
		metric[url.ID] = url.Hits
		return metric, nil

	}
	cursor, err := collection.Find(context.Background(), map[string]string{"username": username})
	if err != nil {
		return nil, err
	}
	var url urlStruct
	metrics := make(map[string]int)
	for cursor.Next(context.Background()) {
		if err = cursor.Decode(&url); err != nil {
			return nil, err
		}
		metrics[url.ID] = url.Hits
	}
	return metrics, nil
}

func (da Adapter) GetURL(id string) (string, error) {
	collection := da.client.Database("urlshortener").Collection("urls")
	var url urlStruct
	err := collection.FindOne(context.Background(), map[string]string{"id": id}).Decode(&url)
	log.Println(id)
	log.Println(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return url.URL, nil
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

type urlStruct struct {
	ID       string `bson:"id"`
	URL      string `bson:"url"`
	Username string `bson:"username"`
	Hits     int    `bson:"hits"`
}
