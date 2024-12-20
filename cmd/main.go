package main

import (
	"fmt"
	"log"
	"os"

	sentry "github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/handler"
)

func main() {

	err := serve()
	if err != nil {
		log.Fatal(err)
	}
}

func serve() error {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           os.Getenv("SENTRY_DSN"),
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate:   1.0,
		ProfilesSampleRate: 1.0,
		Debug:              true,
		Environment:        "debug",
	}); err != nil {
		return fmt.Errorf("sentry initialization failed: %v", err)
	}

	r := gin.Default()
	r.Use(sentrygin.New(sentrygin.Options{Repanic: true}))

	h := handler.ContestHandler{}

	r.GET("/", h.HandleGet)
	r.POST("/", h.HandlePost)
	r.GET("/success", h.HandlePostSuccess)

	err := r.Run(":8080")
	return err
}
