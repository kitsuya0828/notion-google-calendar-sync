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

			if event.NotionEventID == "" { // Not yet added to Notion
				notionEventID, err := notioncalendar.CreateEvent(ctx, notionClient, cfg.NotionDatabaseID, event)
				if err != nil {
					return err
				}
				event.NotionEventID = notionEventID
				err = googlecalendar.UpdateEvent(googleCalendarService, cfg.GoogleCalendarID, event)
				if err != nil {
					return err
				}
			} else if event.GoogleCalendarEventID == "" { // // Not yet added to Google Calendar
				googleCalendarEventID, err := googlecalendar.InsertEvent(googleCalendarService, cfg.GoogleCalendarID, event)
				if err != nil {
					return err
				}
				event.GoogleCalendarEventID = googleCalendarEventID
				err = notioncalendar.UpdateEvent(ctx, notionClient, event)
				if err != nil {
					return err
				}
			}

			err = db.AddEvent(ctx, firestoreClient, event)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
