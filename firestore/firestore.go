package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

const (
	collectionID = "events"
)

func CreateClient(ctx context.Context, projectID string) *firestore.Client {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}

func AddEvent(ctx context.Context, client *firestore.Client, event *Event) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	event.UUID = uuid.String()

	_, err = client.Collection(collectionID).Doc(uuid.String()).Create(ctx, event)
	if err != nil {
		return err
	}
	return nil
}
