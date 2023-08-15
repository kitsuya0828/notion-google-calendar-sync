package run

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/db"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/googlecalendar"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/notioncalendar"
	"github.com/dstotijn/go-notion"
	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
)

func checkAdd(
	ctx context.Context,
	cfg Config,
	notionClient *notion.Client,
	googleCalendarService *calendar.Service,
	notionEvents []*db.Event,
	googleCalendarEvents []*db.Event,
	firestoreClient *firestore.Client,
) error {
	events := append(notionEvents, googleCalendarEvents...)
	for _, event := range events {
		if event.UUID == "" { // Not yet added to the database
			uuid, err := uuid.NewRandom()
			if err != nil {
				return err
			}
			event.UUID = uuid.String()

			if event.NotionEventID == "" { // New Google Calendar event
				notionEventID, err := notioncalendar.CreateEvent(ctx, notionClient, cfg.NotionDatabaseID, event)
				if err != nil {
					return err
				}
				event.NotionEventID = notionEventID
			} else if event.GoogleCalendarEventID == "" { // New Notion event
				googleCalendarEventID, err := googlecalendar.InsertEvent(googleCalendarService, cfg.GoogleCalendarID, event)
				if err != nil {
					return err
				}
				event.GoogleCalendarEventID = googleCalendarEventID
			}

			err = db.AddEvent(ctx, firestoreClient, event)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
