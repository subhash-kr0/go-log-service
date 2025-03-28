package main

import (
	"net/http"
	"os"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func initLogger() {
	// Log as JSON
	log.SetFormatter(&logrus.JSONFormatter{})
	
	// Output to stdout and file
	file, err := os.OpenFile("/var/log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	
	// Set log level
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	initLogger()
	
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Handling info request")
		w.Write([]byte("Logged info"))
	})
	
	http.HandleFunc("/warn", func(w http.ResponseWriter, r *http.Request) {
		log.Warn("Handling warn request")
		w.Write([]byte("Logged warn"))
	})
	
	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Handling debug request")
		w.Write([]byte("Logged debug"))
	})
	
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		log.Error("Handling error request")
		w.Write([]byte("Logged error"))
	})
	
	log.Info("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}