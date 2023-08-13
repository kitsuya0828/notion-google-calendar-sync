package run

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Kitsuya0828/notion-googlecalendar-sync/firestore"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/googlecalendar"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/notion"
	"github.com/caarlos0/env/v9"
)

type Config struct {
	NotionToken      string `env:"NOTION_TOKEN,notEmpty"`
	NotionDatabaseID string `env:"NOTION_DATABASE_ID,notEmpty"`
	GoogleCalendarID string `env:"GOOGLE_CALENDAR_ID,notEmpty"`
	ProjectID        string `env:"PROJECT_ID,notEmpty"`
}

func Sync() {
	ctx := context.Background()

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	notionClient := notion.NewClient(cfg.NotionToken)
	notionEvents, err := notion.ListEvents(ctx, notionClient, cfg.NotionDatabaseID)
	if err != nil {
		log.Fatalf("failed to get events from Notion: %v\n", err)
	}

	if result, err := json.Marshal(notionEvents); err == nil {
		fmt.Println(string(result))
	}

	googleCalendarService, err := googlecalendar.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}
	googleCalendarEvents, err := googlecalendar.ListEvents(googleCalendarService, cfg.GoogleCalendarID)
	if err != nil {
		log.Fatalf("failed to get events from Google Calendar: %v\n", err)
	}
	if result, err := json.Marshal(googleCalendarEvents); err == nil {
		fmt.Println(string(result))
	}

	client := firestore.CreateClient(ctx, cfg.ProjectID)
	firestore.InsertEvents(ctx, client, googleCalendarEvents)
}
