package gcs_s3_sync

import (
        "context"
        "log"
)

type PubSubMessage struct {
        Data []byte `json:"data"`
}

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(ctx context.Context, m PubSubMessage) error {
        name := string(m.Data) // Automatically decoded from base64.
        if name == "" {
                name = "World"
        }
        log.Printf("Hello, %s!", name)
        return nil
}