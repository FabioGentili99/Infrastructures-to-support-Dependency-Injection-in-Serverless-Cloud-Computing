package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	Id             string `json:"id"`
	ServiceName    string `json:"ServiceName"`
	ServiceAddress string `json:"ServiceAddress"`
	Admin          string `json:"Admin"`
	Password       string `json:"Password"`
	Bucket         string `json:"Bucket"`
}

type Payload struct {
	Message string `json:"message"`
}

var logger = logrus.New()

var injectorURL string

var ids []string

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

func handler(w http.ResponseWriter, r *http.Request) {
	var p Payload
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	tsMillis, err := strconv.ParseInt(p.Message, 10, 64)
	if err != nil {
		http.Error(w, "Invalid timestamp", http.StatusBadRequest)
		return
	}

	start := time.Now()

	//resp, err := http.Get("http://injector.default.svc.cluster.local/services/hello")
	resp, err := http.Get(injectorURL + "/services/minio")
	if err != nil {
		http.Error(w, "Failed to reach injector", 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		http.Error(w, "Service not found", 404)
		return
	}
	end := time.Now()
	logger.Infof("Service retrieved in %.3f ms", float64(end.Sub(start).Nanoseconds())/1e6)

	var svc Service
	if err := json.NewDecoder(resp.Body).Decode(&svc); err != nil {
		logger.Error(w, "Invalid response", 500) //////////////////
		return
	}

	start = time.Now()

	minio, err := NewMinio(svc.ServiceAddress, svc.Admin, svc.Password, svc.Bucket)
	if err != nil {
		logger.Error(w, "Failed to create MinIO client: "+err.Error(), http.StatusInternalServerError) /////////////////
		return
	}
	err = minio.Upload(context.Background(), time.Now().UTC().String()+".txt", []byte("Hello from Go"), "text/plain")
	if err != nil {
		logger.Fatal("Upload failed:", err) /////////////////////
	}

	end = time.Now()
	logger.Infof("Service invoked in %.3f ms", float64(end.Sub(start).Nanoseconds())/1e6)

	finish := time.Now().UnixMilli()
	total_latency := finish - tsMillis
	logger.Infof("Total latency is %s ms", strconv.FormatInt(total_latency, 10))

	//body, _ := io.ReadAll(targetResp.Body)
	w.Write([]byte("File uploaded successfully to MinIO!"))
}

func main() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)    // Log level
	logger.SetFormatter(&CSVFormatter{}) // Use custom CSV formatter

	// Get env vars
	injectorURL = os.Getenv("INJECTOR_URL")
	if injectorURL == "" {
		injectorURL = "http://injector.default.svc.cluster.local"
	}

	http.HandleFunc("/", handler)
	logger.Infof("Function invoker running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
