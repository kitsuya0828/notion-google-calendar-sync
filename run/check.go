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

func checkUpdate(
	ctx context.Context,
	cfg Config,
	notionClient *notion.Client,
	googleCalendarService *calendar.Service,
	notionEvents []*db.Event,
	googleCalendarEvents []*db.Event,
	firestoreClient *firestore.Client,
) error {
	events, err := db.ListEvents(ctx, firestoreClient)
	if err != nil {
		return err
	}
	for _, event := range events {
		isNotionDeleted := false
		notionEventsIDMap := getEventsIDMap(notionEvents)
		// Check if the event has been deleted on Notion
		notionEvent, ok := notionEventsIDMap[event.UUID]
		if !ok {
			isNotionDeleted = true
		}

		isGoogleCalendarDeleted := false
		googleCalendarEventsIDMap := getEventsIDMap(googleCalendarEvents)
		// Check if the event has been deleted on Google Calendar
		googelCalendarEvent, ok := googleCalendarEventsIDMap[event.UUID]
		if !ok {
			isGoogleCalendarDeleted = true
		}

		// If the event is deleted either on Notion or Google Calendar
		// TODO: Maintain consistency of events
		if isNotionDeleted && !isGoogleCalendarDeleted {
			err := googlecalendar.DeleteEvent(googleCalendarService, cfg.GoogleCalendarID, event)
			if err != nil {
				return err
			}
			err = db.DeleteEvent(ctx, firestoreClient, event)
			if err != nil {
				return err
			}
			continue
		} else if !isNotionDeleted && isGoogleCalendarDeleted {
			err := notioncalendar.DeleteEvent(ctx, notionClient, event)
			if err != nil {
				return err
			}
			err = db.DeleteEvent(ctx, firestoreClient, event)
			if err != nil {
				return err
			}
			continue
		}

		correctEvent, isNotionUpdated, isGoogleCalendarUpdated := getCorrectEvent(event, notionEvent, googelCalendarEvent)
		if isNotionUpdated || isGoogleCalendarUpdated {
			err := db.SetEvent(ctx, firestoreClient, correctEvent)
			if err != nil {
				return err
			}
			if isNotionUpdated {
				err := googlecalendar.UpdateEvent(googleCalendarService, cfg.GoogleCalendarID, correctEvent)
				if err != nil {
					return err
				}
			}
			if isGoogleCalendarUpdated {
				err := notioncalendar.UpdateEvent(ctx, notionClient, correctEvent)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
