package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI = os.Getenv("MONGO_URI")
var cache = make(map[string]Service)

type Injector struct {
	logger         *logrus.Logger
	dbUrl          string
	dbName         string
	collectionName string
	client         *mongo.Client
	collection     *mongo.Collection
}

// Custom CSV Formatter
type CSVFormatter struct{}

func (f *CSVFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Convert timestamp to ISO 8601 format
	timestamp := entry.Time.UTC().Format("2006-01-02T15:04:05.000Z")
	// Format the log as CSV: timestamp,level,logger,message
	logMsg := fmt.Sprintf("%s,%s,%s,%s\n",
		timestamp, "CALLER", entry.Level.String(), entry.Message)
	return []byte(logMsg), nil
}

func NewInjector() *Injector {
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	injector := &Injector{
		logger:         logger,
		dbUrl:          mongoURI,
		dbName:         "services",
		collectionName: "services",
	}

	injector.connect()
	return injector
}

func (i *Injector) connect() {
	clientOptions := options.Client().ApplyURI(i.dbUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	i.client = client
	i.collection = i.client.Database(i.dbName).Collection(i.collectionName)
}

type Service struct {
	ID             string `bson:"id"`
	ServiceName    string `bson:"ServiceName"`
	ServiceAddress string `bson:"ServiceAddress"`
}

func (i *Injector) RegisterService(id, name, address string) error {
	service := Service{
		ID:             id,
		ServiceName:    name,
		ServiceAddress: address,
	}
	_, err := i.collection.InsertOne(context.TODO(), service)
	if err != nil {
		return err
	}
	fmt.Println("1 document inserted")
	return nil
}

func (i *Injector) GetServiceById(id string) (Service, error) {

	// Check if the service is in cache
	if service, found := cache[id]; found {
		return service, nil
	}

	var service Service
	err := i.collection.FindOne(context.TODO(), bson.D{{Key: "id", Value: id}}).Decode(&service)
	if err != nil {
		return Service{}, err
	}

	// Store in cache
	cache[id] = service
	return service, nil
}
