package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Payload struct {
	Message string `json:"message"`
}

var logger = logrus.New()

func main() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)    // Log level
	logger.SetFormatter(&CSVFormatter{}) // Use custom CSV formatter

	inj := NewInjector()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

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

		svc, err := inj.GetServiceById("hello")

		if err != nil {
			http.Error(w, "Service not found", 404)
			return
		}
		end := time.Now()
		logger.Infof("Service retrieved in %.3f ms", float64(end.Sub(start).Nanoseconds())/1e6)

		start = time.Now()

		// Call the discovered function
		targetResp, err := invoke(svc.ServiceAddress)
		if err != nil {
			http.Error(w, "Failed to call target", 500)
			return
		}
		end = time.Now()
		logger.Infof("Service invoked in %.3f ms", float64(end.Sub(start).Nanoseconds())/1e6)

		finish := time.Now().UnixMilli()
		total_latency := finish - tsMillis
		logger.Infof("Total latency is %s ms", strconv.FormatInt(total_latency, 10))

		w.Write([]byte("Response from hello-world:\n"))
		w.Write([]byte(targetResp))
	})
	logger.Infof("Function invoker running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func invoke(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
