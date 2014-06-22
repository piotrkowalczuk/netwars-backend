package middleware

import (
    "github.com/piotrkowalczuk/netwars-backend/service"
    "github.com/go-martini/martini"
    "net/http"
)

func Sentry(config service.SentryConfig) martini.Handler {
    return func(c martini.Context, req *http.Request) {
        sentry := service.NewSentry(config, req);

        c.Map(sentry)
    }
}
