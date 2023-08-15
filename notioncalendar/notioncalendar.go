package notioncalendar

import (
	"context"
	"strings"
	"time"

	"github.com/Kitsuya0828/notion-googlecalendar-sync/db"
	"github.com/dstotijn/go-notion"
)

func NewClient(token string) *notion.Client {
	client := notion.NewClient(token)
	return client
}

func ListEvents(ctx context.Context, client *notion.Client, databaseID string, tz string) ([]*db.Event, error) {
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

	events := []*db.Event{}

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
			event := &db.Event{NotionEventID: page.ID}

			props, ok := page.Properties.(notion.DatabasePageProperties)
			if !ok {
				continue
			}

			for key, prop := range props {
				switch pt := prop.Type; pt {
				case "title":
					titles := make([]string, len(prop.Title))
					for _, rt := range prop.Title {
						titles = append(titles, rt.PlainText)
					}
					event.Title = strings.Join(titles, "\n")
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
					descriptions := make([]string, len(prop.RichText))
					for _, rt := range prop.RichText {
						if key == "UUID" {
							event.UUID = rt.PlainText
							break
						} else {
							descriptions = append(descriptions, rt.PlainText)
						}
					}
					if len(descriptions) > 0 {
						event.Description = strings.Join(descriptions, "\n")
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
							event.EndTime = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, 1)
						}
					} else {
						if event.IsAllday { // All day (1 day)
							event.EndTime = event.StartTime.AddDate(0, 0, 1)
						} else {
							// If no end time is specified and it is not an all day event, set the duration to 1 hour
							event.EndTime = event.StartTime.Add(time.Hour)
						}
					}
					// if prop.Date.TimeZone != nil {
					// 	fmt.Println(prop.Date.TimeZone)
					// }
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

func CreateEvent(ctx context.Context, client *notion.Client, databaseID string, event *db.Event) (string, error) {
	date := &notion.Date{
		Start: notion.NewDateTime(event.StartTime, !event.IsAllday),
	}

	if event.IsAllday { // All day event
		endTime := notion.NewDateTime(event.EndTime.AddDate(0, 0, -1), false)
		if date.Start != endTime { // All day event (more than 2 days)
			date.End = &endTime
		}
	} else {
		endTime := notion.NewDateTime(event.EndTime, true)
		date.End = &endTime
	}

	params := notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   databaseID,
		DatabasePageProperties: &notion.DatabasePageProperties{
			"title": notion.DatabasePageProperty{
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: event.Title,
						},
					},
				},
			},
			"説明": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: event.Description,
						},
					},
				},
			},
			"Date": notion.DatabasePageProperty{
				Date: date,
			},
			"UUID": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: event.UUID,
						},
					},
				},
			},
		},
	}

	page, err := client.CreatePage(ctx, params)
	if err != nil {
		return "", err
	}
	return page.ID, nil
}

func UpdateEvent(ctx context.Context, client *notion.Client, event *db.Event) error {
	date := &notion.Date{
		Start: notion.NewDateTime(event.StartTime, !event.IsAllday),
	}

	if event.IsAllday { // All day event
		endTime := notion.NewDateTime(event.EndTime.AddDate(0, 0, -1), false)
		if date.Start != endTime { // All day event (more than 2 days)
			date.End = &endTime
		}
	} else {
		endTime := notion.NewDateTime(event.EndTime, true)
		date.End = &endTime
	}

	params := notion.UpdatePageParams{
		DatabasePageProperties: notion.DatabasePageProperties{
			"title": notion.DatabasePageProperty{
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: event.Title,
						},
					},
				},
			},
			"説明": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: event.Description,
						},
					},
				},
			},
			"Date": notion.DatabasePageProperty{
				Date: date,
			},
			"UUID": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: event.UUID,
						},
					},
				},
			},
		},
	}

	_, err := client.UpdatePage(ctx, event.NotionEventID, params)
	if err != nil {
		return err
	}
	return nil
}
