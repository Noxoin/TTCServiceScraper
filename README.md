# TTC Service Scraper

This is a script that will scrape https://ttc.ca/https://www.ttc.ca/service-advisories/subway-service and populate a Google Calendar with the appropriate events.

My version of the Google Calendar: https://calendar.google.com/calendar/u/0?cid=MG1yYjIza3M1NTNpYjhkZjgyYjN2Y2pwazRAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ

## How to run

1. Use [oauth2l](https://github.com/google/oauth2l) to fetch the AUTHORIZATION token.
2. Set the environment variables found in calendar.go
3. Run with `go run ./...`