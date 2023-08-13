package event

import "time"

// Event represents an event to be stored in the database
type Event struct {
	Title            string    `firestore:"title"`
	StartTime        time.Time `firestore:"start_time"`
	EndTime          time.Time `firestore:"end_time"`
	CreatedTime      time.Time `firestore:"created_time"`
	UpdatedTime      time.Time `firestore:"updated_time"`
	Tags             []string  `firestore:"tags"`
	NotionID         string    `firestore:"notion_id"`
	GoogleCalendarID string    `firestore:"google_calendar_id"`
	Description      string    `firestore:"description"`
}
