package run

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
	events, err := notion.GetEvents(ctx, notionClient, cfg.NotionDatabaseID)
	if err != nil {
		log.Fatalf("failed to get events from Notion: %v\n", err)
	}

	if result, err := json.Marshal(events); err == nil {
		fmt.Println(string(result))
	}
}
