package main

import (
	"context"
	"fmt"
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/jomei/notionapi"
)

type Config struct {
	NotionToken     string `env:"NOTION_TOKEN,notEmpty"`
	NotionDatabaseID string `env:"NOTION_DATABASE_ID,notEmpty"`
}

type Data struct {
	Name notionapi.TitleProperty `json:"title"`
}

func main() {
	ctx := context.Background()

	// Parse Slack environment variables
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	client := notionapi.NewClient(notionapi.Token(cfg.NotionToken))

	req := &notionapi.DatabaseQueryRequest{}
	for {
		response, err := client.Database.Query(ctx, notionapi.DatabaseID(cfg.NotionDatabaseID), req)
		if err != nil {
			log.Fatal(err)
		}

		results := response.Results
		for _, r := range results {
			for _, property := range r.Properties {
				switch property.GetType() {
				case "title":
					p, ok := property.(*notionapi.TitleProperty)	// Type Assertion
					if ok {
						for _, t := range p.Title {
							fmt.Println(t.PlainText)
						}
					} else {
						log.Fatal("error: Type Assertion")
					}
				}
			}
		}

		if response.HasMore {
			req.StartCursor = response.NextCursor
		} else {
			break
		}
	}
}