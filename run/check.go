package run

import (
	"context"

	"cloud.google.com/go/firestore"
	fs "github.com/Kitsuya0828/notion-googlecalendar-sync/firestore"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/googlecalendar"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/notion"
	"github.com/google/uuid"
	"github.com/jomei/notionapi"
	"google.golang.org/api/calendar/v3"
)

func checkAdd(
	ctx context.Context,
	cfg Config,
	notionClient *notionapi.Client,
	googleCalendarService *calendar.Service,
	notionEvents []*fs.Event,
	googleCalendarEvents []*fs.Event,
	firestoreClient *firestore.Client,
) error {
	events := append(notionEvents, googleCalendarEvents...)
	for _, event := range events {
		if event.UUID == "" {
			uuid, err := uuid.NewRandom()
			if err != nil {
				return err
			}
			event.UUID = uuid.String()

			if event.NotionEventID == "" {	// New Google Calendar event 
				notionEventID, err := notion.CreateEvent(ctx, notionClient, cfg.NotionDatabaseID, event)
				if err != nil {
					return err
				}
				event.NotionEventID = notionEventID

			} else if event.GoogleCalendarEventID == "" {	// New Notion event
				googleCalendarEventID, err := googlecalendar.InsertEvent(googleCalendarService, cfg.GoogleCalendarID, event)
				if err != nil {
					return err
				}
				event.GoogleCalendarEventID = googleCalendarEventID
			}

			err = fs.AddEvent(ctx, firestoreClient, event)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
