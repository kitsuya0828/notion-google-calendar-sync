package db

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
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
	uuid := event.UUID
	if uuid == "" {
		return fmt.Errorf("no UUID is set for the event: %v", event)
	}

	_, err := client.Collection(collectionID).Doc(uuid).Create(ctx, event)
	if err != nil {
		return err
	}
	return nil
}

func UpdateEvent(ctx context.Context, client *firestore.Client, event *Event) error {
	_, err := client.Collection(collectionID).Doc(event.UUID).Set(ctx, event)
	if err != nil {
		return err
	}
	return nil
}

func DeleteEvent(ctx context.Context, client *firestore.Client, event *Event) error {
	_, err := client.Collection(collectionID).Doc(event.UUID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func ListEvents(ctx context.Context, client *firestore.Client) ([]*Event, error) {
	iter := client.Collection(collectionID).Documents(ctx)
	events := []*Event{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		// fmt.Println(doc.Data())

		var event Event
		err = doc.DataTo(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}
