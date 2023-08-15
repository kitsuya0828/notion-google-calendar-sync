package run

import (
	"context"
	"fmt"
	"log"

	"github.com/Kitsuya0828/notion-googlecalendar-sync/firestore"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/googlecalendar"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/notion"
	"github.com/caarlos0/env/v9"
)

type Config struct {
	NotionToken           string `env:"NOTION_TOKEN,notEmpty"`
	NotionDefaultTimeZone string `env:"NOTION_DEFAULT_TIMEZONE,notEmpty"`
	NotionDatabaseID      string `env:"NOTION_DATABASE_ID,notEmpty"`
	GoogleCalendarID      string `env:"GOOGLE_CALENDAR_ID,notEmpty"`
	ProjectID             string `env:"PROJECT_ID,notEmpty"`
}

func Run() {
	ctx := context.Background()

	// Parse environment varibles
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	// Get future events in Notion database
	notionClient := notion.NewClient(cfg.NotionToken)
	notionEvents, err := notion.ListEvents(ctx, notionClient, cfg.NotionDatabaseID, cfg.NotionDefaultTimeZone)
	if err != nil {
		log.Fatalf("failed to get events from Notion: %v\n", err)
	}

	// Get future events in Google Calendar
	googleCalendarService, err := googlecalendar.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}
	googleCalendarEvents, err := googlecalendar.ListEvents(googleCalendarService, cfg.GoogleCalendarID)
	if err != nil {
		log.Fatalf("failed to get events from Google Calendar: %v\n", err)
	}

	// Initialize Firestore client
	firestoreClient := firestore.CreateClient(ctx, cfg.ProjectID)
	defer firestoreClient.Close()

	// Check if new events have been added
	// err = checkAdd(ctx, cfg, notionClient, googleCalendarService, notionEvents, googleCalendarEvents, firestoreClient)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	for _, event := range notionEvents {
		// firestore.AddEvent(ctx, client, event)
		fmt.Println(event)
	}
	for _, event := range googleCalendarEvents {
		// firestore.AddEvent(ctx, client, event)
		fmt.Println(event)
	}
}
