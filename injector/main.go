package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sirupsen/logrus"
)

type Service struct {
	id             string `json:"id" bson:"id"`
	ServiceName    string `json:"ServiceName" bson:"ServiceName"`
	ServiceAddress string `json:"ServiceAddress" bson:"ServiceAddress"`
}

var collection *mongo.Collection
var cache sync.Map

var logger = logrus.New()

// Custom CSV Formatter
type CSVFormatter struct{}

func (f *CSVFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Convert timestamp to ISO 8601 format
	timestamp := entry.Time.UTC().Format("2006-01-02T15:04:05.000Z")
	// Format the log as CSV: timestamp,level,logger,message
	logMsg := fmt.Sprintf("%s,%s,%s,%s\n",
		timestamp, "INJECTOR", entry.Level.String(), entry.Message)
	return []byte(logMsg), nil
}

func main() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)    // Log level
	logger.SetFormatter(&CSVFormatter{}) // Use custom CSV formatter

	// Get env vars
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://mongo:27017"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Infof("MongoDB connection error: %v", err)
	}
	collection = client.Database("services").Collection("services")
	logger.Infof("Connected to MongoDB")

	// Gin router
	r := gin.Default()
	r.GET("/services/:id", getServiceHandler)
	r.GET("/health", healthCheckHandler)

	logger.Infof("Injector API running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		logger.Infof("Failed to run server: %v", err)
	}
}

func getServiceHandler(c *gin.Context) {
	id := c.Param("id")
	logger.Infof("Fetching service with ID: %s", id)

	start := time.Now()

	// Check if the service is in cache
	if val, ok := cache.Load(id); ok {
		end := time.Now()
		logger.Infof("Service retrieved in %.3f ms", float64(end.Sub(start).Nanoseconds())/1e6)
		c.JSON(200, val)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var service Service

	// Find the service in MongoDB
	err := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&service)
	if err != nil {
		logger.Infof("Error finding service with id '%s': %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}

	end := time.Now()
	logger.Infof("Service retrieved in %.3f ms", float64(end.Sub(start).Nanoseconds())/1e6)
	// Store in cache
	cache.Store(id, service)

	c.JSON(http.StatusOK, service)
}

func healthCheckHandler(c *gin.Context) {
	logger.Infof("Health check endpoint hit")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
