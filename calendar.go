package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var CALENDAR_ID = "0mrb23ks553ib8df82b3vcjpk4@group.calendar.google.com"
var SECRETS_FILE = "client_secret.json"

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
	conf, err := google.ConfigFromJSON([]byte(readFile(SECRETS_FILE)), "https://www.googleapis.com/auth/calendar", "https://www.googleapis.com/auth/calendar.events")
	if err != nil {
		log.Fatalf("unable to get oauth config: %v", err)
	}
	state := "helloworldme"
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Visit this URL in your browser:\n\n%s\n\n", url)
	tok := make(chan *oauth2.Token)
	var wg sync.WaitGroup
	wg.Add(1)
	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()

		if s := r.URL.Query().Get("state"); s != state {
			http.Error(w, fmt.Sprintf("Invalid state: %s", s), http.StatusUnauthorized)
			return
		}

		code := r.URL.Query().Get("code")
		token, err := conf.Exchange(ctx, code)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exchange error: %s", err), http.StatusServiceUnavailable)
			return
		}

		tokenJSON, err := json.MarshalIndent(token, "", "  ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Token parse error: %s", err), http.StatusServiceUnavailable)
			return
		}

		w.Write(tokenJSON)
		tok <- token
	})
	server := http.Server{
		Addr: ":8080",
	}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	receivedToken := <-tok
	wg.Wait()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
	client := conf.Client(ctx, receivedToken)
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
	resp, err := c.srv.Events.Insert(CALENDAR_ID, calEvent).Do()
	log.Printf("Added event: %s\n", resp.Summary)
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
	return string(contents)
}
