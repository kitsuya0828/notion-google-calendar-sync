package db

import "time"

// ColorMap is a map to convert from [Notion] colors to [Google Calendar] colors
//
// [Notion]: https://developers.notion.com/reference/property-object#multi-select
// [Google Calendar]: https://developers.google.com/calendar/api/v3/reference/colors/get?hl=ja
var ColorMap = map[string]string{
	"blue":    "9",
	"brown":   "10",
	"default": "1",
	"gray":    "8",
	"green":   "2",
	"orange":  "6",
	"pink":    "4",
	"purple":  "3",
	"red":     "11",
	"yellow":  "5",
}

// Event represents an event to be stored in the database
type Event struct {
	UUID                  string    `firestore:"uuid"`
	Title                 string    `firestore:"title"`
	StartTime             time.Time `firestore:"start_time"`
	EndTime               time.Time `firestore:"end_time"`
	CreatedTime           time.Time `firestore:"created_time"`
	UpdatedTime           time.Time `firestore:"updated_time"`
	Color                 string    `firestore:"color"`
	IsAllday              bool      `firestore:"is_all_day"`
	NotionEventID         string    `firestore:"notion_event_id"`
	GoogleCalendarEventID string    `firestore:"google_calendar_event_id"`
	Description           string    `firestore:"description"`
}
