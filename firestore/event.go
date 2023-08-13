package firestore

import "time"

// Event represents an event to be stored in the database
type Event struct {
	UUID                  string    `firestore:"uuid"`
	Title                 string    `firestore:"title"`
	StartTime             time.Time `firestore:"start_time"`
	EndTime               time.Time `firestore:"end_time"`
	CreatedTime           time.Time `firestore:"created_time"`
	UpdatedTime           time.Time `firestore:"updated_time"`
	Tags                  []string  `firestore:"tags"`
	NotionEventID         string    `firestore:"notion_event_id"`
	GoogleCalendarEventID string    `firestore:"google_calendar_event_id"`
	Description           string    `firestore:"description"`
}
