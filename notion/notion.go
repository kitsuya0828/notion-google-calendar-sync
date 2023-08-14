package notion

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Kitsuya0828/notion-googlecalendar-sync/firestore"
	"github.com/jomei/notionapi"
)

func NewClient(token string) *notionapi.Client {
	client := notionapi.NewClient(notionapi.Token(token))
	return client
}

func ListEvents(ctx context.Context, client *notionapi.Client, databaseID string) ([]*firestore.Event, error) {
	now := notionapi.Date(time.Now())
	req := &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: "Date",
			Date: &notionapi.DateFilterCondition{
				After: &now,
			},
		},
	}
	events := []*firestore.Event{}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	time.Local = jst

	for {
		response, err := client.Database.Query(ctx, notionapi.DatabaseID(databaseID), req)
		if err != nil {
			return nil, err
		}

		for _, r := range response.Results {
			event := &firestore.Event{NotionEventID: r.ID.String()}
			for label, property := range r.Properties {
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
							event.Color = o.Color.String()
							break
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
							if label == "UUID" {
								event.UUID = rt.Text.Content
								break
							} else {
								event.Description += fmt.Sprintln(rt.Text.Content)
							}
						}
					}
				case "date":
					prop, ok := property.(*notionapi.DateProperty)
					if ok {
						fmt.Println(prop)
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

func CreateEvent(ctx context.Context, client *notionapi.Client, databaseID string, event *firestore.Event) (string, error) {
	startTime := notionapi.Date(event.StartTime)
	endTime := notionapi.Date(event.EndTime)
	if event.EndTime.Hour() == 0 && event.EndTime.Minute() == 0 {	// For convenience of Notion display
		endTime = notionapi.Date(event.EndTime.Add(- time.Second))
	}
	date := &notionapi.DateObject{
		Start: &startTime,
		End: &endTime,
	}


	req := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type: "database_id",
			DatabaseID: notionapi.DatabaseID(databaseID),
		},
		Properties: notionapi.Properties{
			"title": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: event.Title,
						},
					},
				},
			},
			"説明": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: event.Description,
						},
					},
				},
			},
			"Date": &notionapi.DateProperty{
				Date: date,
			},
		},
	}
	response, err := client.Page.Create(ctx, req)
	if err != nil {
		return "", err
	}
	return response.ID.String() ,nil
}