package googlecalendar

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Kitsuya0828/notion-google-calendar-sync/event"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func NewService(ctx context.Context) (*calendar.Service, error) {
	return calendar.NewService(ctx, option.WithCredentialsFile("credentials.json"))
}

func GetEvents(service *calendar.Service) ([]*event.Event, error) {
	events, err := service.Events.List("kitsuyaazuma@gmail.com").TimeMin(time.Now().Format(time.RFC3339)).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%s (%s)\n", item.Summary, date)
		}
	}
	return []*event.Event{}, nil
}
