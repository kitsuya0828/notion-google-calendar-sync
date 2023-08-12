package run

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Kitsuya0828/notion-google-calendar-sync/googlecalendar"
	"github.com/Kitsuya0828/notion-google-calendar-sync/notion"
	"github.com/caarlos0/env/v9"
)

type Config struct {
	NotionToken      string `env:"NOTION_TOKEN,notEmpty"`
	NotionDatabaseID string `env:"NOTION_DATABASE_ID,notEmpty"`
}

func Sync() {
	ctx := context.Background()

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	notionClient := notion.NewClient(cfg.NotionToken)
	notionEvents, err := notion.GetEvents(ctx, notionClient, cfg.NotionDatabaseID)
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
	GoogleCalendarEvents, err := googlecalendar.GetEvents(googleCalendarService)
	if err != nil {
		log.Fatalf("failed to get events from Google Calendar: %v\n", err)
	}
	if result, err := json.Marshal(GoogleCalendarEvents); err == nil {
		fmt.Println(string(result))
	}
}
