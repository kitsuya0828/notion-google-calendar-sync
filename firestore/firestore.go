package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/event"
	"google.golang.org/api/option"
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

func InsertEvents(ctx context.Context, client *firestore.Client, events []*event.Event) error {
	for _, e := range events {
		_, _, err := client.Collection("events").Add(ctx, e)
		if err != nil {
			return err
		}
	}
	return nil
}
