package event

import "time"

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
