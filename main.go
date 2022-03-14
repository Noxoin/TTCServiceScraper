package main

import (
	"context"
	"log"
	"time"
)

type ClosureType string

const (
	Undefined ClosureType = ""
	FullDay   ClosureType = "Full Day"
	NightOnly ClosureType = "Night Only"
)

type Event struct {
	Summary string
	Uri     string
	Type    ClosureType
	Start   time.Time
	End     time.Time
}

func main() {
	ctx := context.Background()
	events, err := runTTCStage(false)
	if err != nil {
		log.Fatalf("error in TTC Stage: %v", err)
	}
	log.Printf("%v", events)
	srv, err := connectToCalendar(ctx)
	if err != nil {
		log.Fatalf("unable to connect to Calendar Service: %v", err)
	}
	c := NewCalendarService(srv, false)
	c.updateCalendar(ctx, events)
}
