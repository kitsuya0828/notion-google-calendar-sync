package googlecalendar

import (
	"context"
	"fmt"
	"time"

	"github.com/Kitsuya0828/notion-googlecalendar-sync/db"
	"golang.org/x/exp/slog"
	"google.golang.org/api/calendar/v3"
)

func NewService(ctx context.Context) (*calendar.Service, error) {
	return calendar.NewService(ctx)
}

func ListEvents(service *calendar.Service, calendarID string) ([]*db.Event, error) {
	events := []*db.Event{}
	result, err := service.Events.List(calendarID).TimeMin(time.Now().Format(time.RFC3339)).Do()
	if err != nil {
		return nil, fmt.Errorf("execute calendar.events.list call: %v", err)
	}

	tz, err := time.LoadLocation(result.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("load location: %v", err)
	}
	time.Local = tz

	for _, item := range result.Items {
		event := &db.Event{
			Title:                 item.Summary,
			GoogleCalendarEventID: item.Id,
			Description:           item.Description,
		}
		slog.Debug("parse item", "item", item)

		createdTime, err := time.Parse(time.RFC3339, item.Created)
		if err != nil {
			return nil, fmt.Errorf("parse created time: %v", err)
		}
		event.CreatedTime = createdTime

		updatedTime, err := time.Parse(time.RFC3339, item.Updated)
		if err != nil {
			return nil, fmt.Errorf("parse updated time: %v", err)
		}
		event.UpdatedTime = updatedTime

		startTime := time.Time{}
		if item.Start.DateTime == "" {
			startTime, err = time.ParseInLocation("2006-01-02", item.Start.Date, tz)
			if err != nil {
				return nil, fmt.Errorf("parse start time: %v", err)
			}
			event.IsAllday = true
		} else {
			startTime, err = time.Parse(time.RFC3339, item.Start.DateTime)
			if err != nil {
				return nil, fmt.Errorf("parse start time: %v", err)
			}
		}
		event.StartTime = startTime

		endTime := time.Time{}
		if item.End.DateTime == "" {
			endTime, err = time.ParseInLocation("2006-01-02", item.End.Date, tz)
			if err != nil {
				return nil, fmt.Errorf("parse end time: %v", err)
			}
		} else {
			endTime, err = time.Parse(time.RFC3339, item.End.DateTime)
			if err != nil {
				return nil, fmt.Errorf("parse end time: %v", err)
			}
		}
		event.EndTime = endTime

		if item.ExtendedProperties != nil {
			uuid, ok := item.ExtendedProperties.Private["uuid"]
			if ok {
				event.UUID = uuid
			}
		}

		for k, v := range db.ColorMap {
			if v == item.ColorId {
				event.Color = k
				break
			}
		}
		events = append(events, event)
	}
	slog.Info("list google calendar events", "num", len(events))
	return events, nil
}

func InsertEvent(service *calendar.Service, calendarID string, event *db.Event) (string, error) {
	startDateTime := &calendar.EventDateTime{
		DateTime: event.StartTime.Format(time.RFC3339),
	}
	endDateTime := &calendar.EventDateTime{
		DateTime: event.EndTime.Format(time.RFC3339),
	}
	if event.IsAllday {
		startDateTime = &calendar.EventDateTime{
			Date: event.StartTime.Format("2006-01-02"),
		}
		endDateTime = &calendar.EventDateTime{
			Date: event.EndTime.Format("2006-01-02"),
		}
	}

	e := &calendar.Event{
		Summary:     event.Title,
		Description: event.Description,
		Start:       startDateTime,
		End:         endDateTime,
		ExtendedProperties: &calendar.EventExtendedProperties{
			Private: map[string]string{
				"uuid": event.UUID,
			},
		},
		ColorId: db.ColorMap[event.Color],
	}

	result, err := service.Events.Insert(calendarID, e).Do()
	if err != nil {
		return "", err
	}
	slog.Info("insert google calendar event", "event", result)
	return result.Id, nil
}

func UpdateEvent(service *calendar.Service, calendarID string, event *db.Event) error {
	startDateTime := &calendar.EventDateTime{
		DateTime: event.StartTime.Format(time.RFC3339),
	}
	endDateTime := &calendar.EventDateTime{
		DateTime: event.EndTime.Format(time.RFC3339),
	}
	if event.IsAllday {
		startDateTime = &calendar.EventDateTime{
			Date: event.StartTime.Format("2006-01-02"),
		}
		endDateTime = &calendar.EventDateTime{
			Date: event.EndTime.Format("2006-01-02"),
		}
	}

	e := &calendar.Event{
		Summary:     event.Title,
		Description: event.Description,
		Start:       startDateTime,
		End:         endDateTime,
		ExtendedProperties: &calendar.EventExtendedProperties{
			Private: map[string]string{
				"uuid": event.UUID,
			},
		},
	}
	if event.Color != "" {
		e.ColorId = db.ColorMap[event.Color]
	}

	result, err := service.Events.Update(calendarID, event.GoogleCalendarEventID, e).Do()
	if err != nil {
		return err
	}
	slog.Info("update google calendar event", "event", result)
	return nil
}

func DeleteEvent(service *calendar.Service, calendarID string, event *db.Event) error {
	err := service.Events.Delete(calendarID, event.GoogleCalendarEventID).Do()
	if err != nil {
		return err
	}
	slog.Info("delete google calendar event", "id", event.GoogleCalendarEventID)
	return nil
}
