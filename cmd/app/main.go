package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func initLogger(logFile *os.File) {
	// Log as JSON
	log.SetFormatter(&logrus.JSONFormatter{})

	// Write logs to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// Set log level
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	// ✅ Ensure logs directory exists
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		fmt.Println("Failed to create logs directory:", err)
		os.Exit(1) // Stop execution if log directory can't be created
	}

	// ✅ Create/Open log file
	logFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file, using default stderr:", err)
		logFile = os.Stdout // Fallback to stdout if log file fails
	}
	defer logFile.Close() // Close file on exit

	initLogger(logFile)

	// ✅ Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Root route hit")
		w.Write([]byte("Welcome to the server!"))
	})

	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Handling info request")
		w.Write([]byte("Logged info"))
	})

	mux.HandleFunc("/warn", func(w http.ResponseWriter, r *http.Request) {
		log.Warn("Handling warn request")
		w.Write([]byte("Logged warn"))
	})

	mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Handling debug request")
		w.Write([]byte("Logged debug"))
	})

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		log.Error("Handling error request")
		w.Write([]byte("Logged error"))
	})

	// ✅ Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Info("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: ", err)
		}
	}()

	// ✅ Wait for termination signal
	<-stop
	log.Info("Shutting down server gracefully...")

	// ✅ Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Info("Server exited properly")
}
