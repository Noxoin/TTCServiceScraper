package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var CALENDAR_ID = os.Getenv("CALENDAR_ID")
var CLIENT_ID = os.Getenv("CLIENT_ID")
var SECRETS_FILE = os.Getenv("SECRETS_FILE")
var AUTHORIZATION = os.Getenv("AUTHORIZATION")

type CalendarService struct {
	srv    *calendar.Service
	dryRun bool
}

func NewCalendarService(srv *calendar.Service, dryRun bool) *CalendarService {
	return &CalendarService{
		srv:    srv,
		dryRun: dryRun,
	}
}

func connectToCalendar(ctx context.Context) (*calendar.Service, error) {
	conf := &oauth2.Config{
		ClientID:     CLIENT_ID,
		ClientSecret: readFile(SECRETS_FILE),
		Scopes:       []string{"https://www.googleapis.com/auth/calendar", "https://www.googleapis.com/auth/calendar.events"},
		Endpoint:     google.Endpoint,
	}
	client := conf.Client(ctx, &oauth2.Token{
		AccessToken: AUTHORIZATION,
	})
	return calendar.NewService(ctx, option.WithHTTPClient(client))
}

func (c *CalendarService) insertCalendarEvent(ctx context.Context, event Event) error {
	calEvent := &calendar.Event{
		Summary:     event.Summary,
		Description: event.Uri,
	}
	if event.Type == Undefined {
		return fmt.Errorf("unable to handle undefined event type: %v", event)
	} else if event.Type == FullDay {
		calEvent.Start = &calendar.EventDateTime{
			Date: event.Start.Format("2006-01-02"),
		}
		calEvent.End = &calendar.EventDateTime{
			Date: event.End.AddDate(0, 0, 1).Format("2006-01-02"),
		}
	} else if event.Type == NightOnly {
		calEvent.Start = &calendar.EventDateTime{
			DateTime: fmt.Sprintf("%sT23:00:00", event.Start.Format("2006-01-02")),
			TimeZone: "America/Toronto",
		}
		calEvent.End = &calendar.EventDateTime{
			DateTime: fmt.Sprintf("%sT06:00:00", event.Start.AddDate(0, 0, 1).Format("2006-01-02")),
			TimeZone: "America/Toronto",
		}
		count := int(event.End.Sub(event.Start).Hours())/int(24) + 1 // occurrance count includes the start and end days
		calEvent.Recurrence = []string{fmt.Sprintf("RRULE:FREQ=DAILY;COUNT=%d", count)}
	}
	if c.dryRun {
		log.Printf("Will send request: %v", calEvent)
		return nil
	}
	_, err := c.srv.Events.Insert(CALENDAR_ID, calEvent).Do()
	return err
}

func (c *CalendarService) insertCalendarEvents(ctx context.Context, events []Event) {
	for _, e := range events {
		err := c.insertCalendarEvent(ctx, e)
		if err != nil {
			log.Printf("%v", err)
		}
	}
}

func (c *CalendarService) fetchExistingEvents(ctx context.Context, minTime time.Time, maxTime time.Time) (*map[string]bool, error) {
	var urls *map[string]bool = &map[string]bool{}
	events, err :=
		c.srv.Events.List(CALENDAR_ID).TimeMax(fmt.Sprintf("%sT00:00:00Z", maxTime.Format("2006-01-02"))).TimeMin(fmt.Sprintf("%sT00:00:00Z", minTime.Format("2006-01-02"))).Do()
	if err != nil {
		return urls, err
	}
	for _, e := range events.Items {
		(*urls)[e.Description] = true
	}
	return urls, nil
}

func (c *CalendarService) updateCalendar(ctx context.Context, events []Event) {
	var minTime, maxTime time.Time
	for _, e := range events {
		if minTime.After(e.Start) {
			minTime = e.Start
		}
		if maxTime.Before(e.End) {
			maxTime = e.End
		}
	}
	maxTime = maxTime.AddDate(0, 0, 1) // Add 1 day to have exclusive end.
	existingEventIds, err := c.fetchExistingEvents(ctx, minTime, maxTime)
	if err != nil {
		log.Printf("unable to get existing events: %v", err)
		return
	}
	var eventsToUpdate []Event
	for _, e := range events {
		if !(*existingEventIds)[e.Uri] {
			eventsToUpdate = append(eventsToUpdate, e)
		}
	}
	c.insertCalendarEvents(ctx, eventsToUpdate)
}

func readFile(filename string) string {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading %q: %v", filename, err)
	}
	return strings.TrimSpace(string(contents))
}
