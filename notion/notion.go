package notion

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Kitsuya0828/notion-googlecalendar-sync/firestore"
	"github.com/dstotijn/go-notion"
	"github.com/jomei/notionapi"
)

func NewClient(token string) *notion.Client {
	client := notion.NewClient(token)
	return client
}

func ListEvents(ctx context.Context, client *notion.Client, databaseID string, tz string) ([]*firestore.Event, error) {
	now := time.Now()
	req := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "Date",
			DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
				Date: &notion.DatePropertyFilter{
					After: &now,
				},
			},
		},
	}

	events := []*firestore.Event{}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}
	time.Local = loc

	for {
		response, err := client.QueryDatabase(ctx, databaseID, req)
		if err != nil {
			return nil, err
		}
		result := response.Results

		for _, page := range result {
			event := &firestore.Event{NotionEventID: page.ID}

			props, ok := page.Properties.(notion.DatabasePageProperties)
			if !ok {
				continue
			}

			for key, prop := range props {
				switch pt := prop.Type; pt {
				case "title":
					for _, rt := range prop.Title {
						event.Title += fmt.Sprintln(rt.PlainText)
					}
				case "multi_select":
					for _, o := range prop.MultiSelect {
						event.Color = string(o.Color)
						break
					}
				case "created_time":
					event.CreatedTime = *prop.CreatedTime
				case "last_edited_time":
					event.UpdatedTime = *prop.LastEditedTime
				case "rich_text":
					for _, rt := range prop.RichText {
						if key == "UUID" {
							event.UUID = rt.PlainText
							break
						} else {
							event.Description += fmt.Sprintln(rt.PlainText)
						}
					}
				case "date":
					event.StartTime = prop.Date.Start.Time
					if !prop.Date.Start.HasTime() { // All day
						st := event.StartTime
						event.StartTime = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, loc)
						event.IsAllday = true
					}
					if prop.Date.End != nil {
						event.EndTime = prop.Date.End.Time
						if !prop.Date.End.HasTime() { // All day (more than 2 days)
							et := event.EndTime
							event.EndTime = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, loc).Add(24 * time.Hour)
						}
					} else {
						if event.IsAllday { // All day (1 day)
							event.EndTime = event.StartTime.Add(24 * time.Hour)
						} else {
							// If no end time is specified and it is not an all day event, set the duration to 1 hour
							event.EndTime = event.StartTime.Add(time.Hour)
						}
					}
					// if prop.Date.TimeZone != nil {
					// 	fmt.Println(prop.Date.TimeZone)
					// }
				default:
					log.Printf("property is not supported: %s\n", pt)
				}
			}
			events = append(events, event)
		}

		if response.HasMore {
			req.StartCursor = *response.NextCursor
		} else {
			break
		}
	}
	return events, nil
}

func CreateEvent(ctx context.Context, client *notionapi.Client, databaseID string, event *firestore.Event) (string, error) {
	startTime := notionapi.Date(event.StartTime)
	endTime := notionapi.Date(event.EndTime)
	// This SDK cannot handle all day events
	if event.EndTime.Hour() == 0 && event.EndTime.Minute() == 0 { // For convenience of Notion display
		endTime = notionapi.Date(event.EndTime.Add(-time.Second))
	}
	date := &notionapi.DateObject{
		Start: &startTime,
		End:   &endTime,
	}

	req := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       "database_id",
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
	return response.ID.String(), nil
}
