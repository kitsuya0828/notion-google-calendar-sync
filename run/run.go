package run

import (
	"context"
	"fmt"

	"github.com/Kitsuya0828/notion-googlecalendar-sync/db"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/googlecalendar"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/notioncalendar"
	"github.com/caarlos0/env/v9"
	"golang.org/x/exp/slog"
)

type Config struct {
	NotionToken           string `env:"NOTION_TOKEN,notEmpty"`
	NotionDefaultTimeZone string `env:"NOTION_DEFAULT_TIMEZONE,notEmpty"`
	NotionDatabaseID      string `env:"NOTION_DATABASE_ID,notEmpty"`
	GoogleCalendarID      string `env:"GOOGLE_CALENDAR_ID,notEmpty"`
	ProjectID             string `env:"PROJECT_ID,notEmpty"`
}

func Run(logger *slog.Logger) error {
	ctx := context.Background()

	// Parse environment varibles
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("parse env: %v", err)
	}

	// List future events in Notion database
	logger.Debug("list notion events")
	notionClient := notioncalendar.NewClient(cfg.NotionToken)
	notionEvents, err := notioncalendar.ListEvents(ctx, notionClient, cfg.NotionDatabaseID, cfg.NotionDefaultTimeZone)
	if err != nil {
		return fmt.Errorf("list notion events: %v", err)
	}

	// List future events in Google Calendar
	logger.Debug("list google calendar events")
	googleCalendarService, err := googlecalendar.NewService(ctx)
	if err != nil {
		return fmt.Errorf("initialize google calendar service: %v", err)
	}
	googleCalendarEvents, err := googlecalendar.ListEvents(logger, googleCalendarService, cfg.GoogleCalendarID)
	if err != nil {
		return fmt.Errorf("list google calendar events: %v", err)
	}

	// Initialize Firestore client
	logger.Debug("initialize firestore client")
	firestoreClient := db.CreateClient(ctx, cfg.ProjectID)
	defer firestoreClient.Close()

	// Check if new events have been added
	logger.Debug("check for added events")
	err = checkAdd(ctx, cfg, notionClient, googleCalendarService, notionEvents, googleCalendarEvents, firestoreClient)
	if err != nil {
		return fmt.Errorf("check for added events: %v", err)
	}

	// List future events again
	logger.Debug("list notion events again")
	notionEvents, err = notioncalendar.ListEvents(ctx, notionClient, cfg.NotionDatabaseID, cfg.NotionDefaultTimeZone)
	if err != nil {
		return fmt.Errorf("list notion events again: %v", err)
	}
	logger.Debug("list google calendar events again")
	googleCalendarEvents, err = googlecalendar.ListEvents(logger, googleCalendarService, cfg.GoogleCalendarID)
	if err != nil {
		return fmt.Errorf("list google calendar events again: %v", err)
	}

	// Check if events have been updated or deleted
	err = checkUpdate(ctx, cfg, notionClient, googleCalendarService, notionEvents, googleCalendarEvents, firestoreClient)
	if err != nil {
		return fmt.Errorf("check for updated or deleted events: %v", err)
	}

	return nil
}
