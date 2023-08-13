package notion

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Kitsuya0828/notion-googlecalendar-sync/event"
	"github.com/jomei/notionapi"
)

func NewClient(token string) *notionapi.Client {
	client := notionapi.NewClient(notionapi.Token(token))
	return client
}

func GetEvents(ctx context.Context, client *notionapi.Client, databaseID string) ([]*event.Event, error) {
	req := &notionapi.DatabaseQueryRequest{}
	events := []*event.Event{}
	for {
		response, err := client.Database.Query(ctx, notionapi.DatabaseID(databaseID), req)
		if err != nil {
			return nil, err
		}

		for _, r := range response.Results {
			event := &event.Event{NotionID: r.ID.String()}
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
						} else {
							event.EndTime = event.StartTime.Add(24 * time.Hour)
						}
					}
				default:
					log.Printf("property is not supported: %s\n", pt)
				}
			}
			events = append(events, event)
		}

		if response.HasMore {
			req.StartCursor = response.NextCursor
		} else {
			break
		}
	}
	return events, nil
}
