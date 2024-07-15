package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"memrizr/account/observability/logger"
	"memrizr/account/observability/prometheus"
	"memrizr/account/observability/tracing"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting server...")

	// initialize data sources
	ds, err := initDS()

	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}

	router, err := inject(ds)

	if err != nil {
		log.Fatalf("Failure to inject data sources: %v\n", err)
	}

	logger.Setup()
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			logger.LogError(err, "Failed to close log file", logrus.Fields{"module": "main", "function": "main"})
		}
	}(logger.LogFile)

	fields := logrus.Fields{"module": "main", "function": "main"}
	logger.LogInfo("Service started", fields)

	err = tracing.InitTracer()
	if err != nil {
		log.Fatalf("failed to initialize OpenTelemetry: %v", err)
	}

	prometheus.InitMetrics()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", srv.Addr)

	// Wait for kill signal of channel
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown data sources
	log.Println("Shutting down database server...")
	if err := ds.close(); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}
	// Shutdown server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
