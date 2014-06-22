package service

import (
	"github.com/getsentry/raven-go"
    "log"
    "net/http"
    "fmt"
)

type Sentry struct {
	client     *raven.Client
    request    *http.Request
}

type SentryConfig struct {
	DSN         string `xml:"dsn"`
}

func NewSentry(config SentryConfig, r *http.Request) *Sentry {
	client, err := raven.NewClient(config.DSN, nil)
	if err != nil {
		log.Fatal(err)
	}

	sentry := Sentry{client, r}

    return &sentry
}

func (s *Sentry) Error(reportedError error) {
	if reportedError == nil {
		return
	}
	
	var err error
    trace := raven.NewStacktrace(0, 2, nil)
    packet := raven.NewPacket(
		reportedError.Error(),
		raven.NewException(reportedError, trace),
		raven.NewHttp(s.request),
	)
    eventID, ch := s.client.Capture(packet, nil)
    if err = <-ch; err != nil {
        log.Fatal(err)
    }
    message := fmt.Sprintf("Captured error with id %s: %q", eventID, reportedError)
    log.Println(message)
}
