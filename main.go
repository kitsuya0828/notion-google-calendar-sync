package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/jomei/notionapi"
)

type Config struct {
	NotionToken      string `env:"NOTION_TOKEN,notEmpty"`
	NotionDatabaseID string `env:"NOTION_DATABASE_ID,notEmpty"`
}

type Event struct {
	Title            string    `json:"title"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	CreatedTime      time.Time `json:"created_time"`
	UpdatedTime      time.Time `json:"updated_time"`
	Tags             []string  `json:"tags"`
	NotionID         string    `json:"notion_id"`
	GoogleCalendarID string    `json:"google_calendar_id"`
	Description      string    `json:"description"`
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
			event := &Event{NotionID: r.ID.String()}
			for _, property := range r.Properties {
				switch pt := property.GetType(); pt {
				case "title":
					prop, ok := property.(*notionapi.TitleProperty)
					if ok {
						for _, t := range prop.Title {
							event.Title = t.PlainText
						}
					}
				case "multi_select":
					prop, ok := property.(*notionapi.MultiSelectProperty)
					if ok {
						for _, o := range prop.MultiSelect {
							event.Tags = append(event.Tags, o.Name)
						}
					}
				case "created_time":
					prop, ok := property.(*notionapi.CreatedTimeProperty)
					if ok {
						event.CreatedTime = prop.CreatedTime
					}
				case "last_edited_time":
					prop, ok := property.(*notionapi.LastEditedTimeProperty)
					if ok {
						event.UpdatedTime = prop.LastEditedTime
					}
				case "rich_text":
					prop, ok := property.(*notionapi.RichTextProperty)
					if ok {
						for _, rt := range prop.RichText {
							event.Description += fmt.Sprintln(rt.Text.Content)
						}
					}
				case "date":
					prop, ok := property.(*notionapi.DateProperty)
					if ok {
						event.StartTime = time.Time(*prop.Date.Start)
						if prop.Date.End != nil {
							event.EndTime = time.Time(*prop.Date.End)
						}
					}
				default:
					log.Printf("property is not supported: %s\n", pt)
				}
			}
			// output, err := json.Marshal(event)
			output, err := json.MarshalIndent(event, "", "\t")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(output))
		}

		if response.HasMore {
			req.StartCursor = response.NextCursor
		} else {
			break
		}
	}
}
