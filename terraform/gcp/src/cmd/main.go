package main

import (
    "os"
	"log"
	"context"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"example.com/gcs_s3_sync"
)

func main() {
	ctx := context.Background()
	if err := funcframework.RegisterEventFunctionContext(ctx, "/", gcs_s3_sync.HelloGCS()); err != nil {
        log.Fatalf("funcframework.RegisterEventFunctionContext: %v\n", err)
	}

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}