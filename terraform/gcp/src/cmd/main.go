package main

import (
    "os"
    "log"
    "context"

    "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
    "cloud.google.com/go/functions/metadata"
    "example.com/gcs_s3_sync"
)

func main() {
    name := "hello_gcs.txt"
    e := gcs_s3_sync.GCSEvent{
        Name: name,
    }
    meta := &metadata.Metadata{
        EventID: "event ID",
    }
    ctx := metadata.NewContext(context.Background(), meta)
    if err := funcframework.RegisterEventFunctionContext(ctx, "/", gcs_s3_sync.HelloGCS(ctx, e)); err != nil {
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