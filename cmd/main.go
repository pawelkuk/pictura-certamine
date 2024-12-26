package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	sentry "github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/config"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/handler"
	contestant "github.com/pawelkuk/pictura-certamine/pkg/domain/contest/repo/contestant"
	entry "github.com/pawelkuk/pictura-certamine/pkg/domain/contest/repo/entry"
	"github.com/pawelkuk/pictura-certamine/pkg/sdk/mail"
	"github.com/pawelkuk/pictura-certamine/pkg/sdk/s3"
)

func main() {

	err := serve()
	if err != nil {
		log.Fatal(err)
	}
}

func serve() error {
	cfg := &config.Config{}
	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("could not parse config: %w", err)
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           cfg.SentryDSN,
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
	db, err := sql.Open("sqlite3", "./data/pictura-certamine.db")
	if err != nil {
		return fmt.Errorf("could not open db: %w", err)
	}
	contestantrepo := &contestant.SQLiteRepo{DB: db}
	entryrepo := &entry.SQLiteRepo{DB: db}
	s3Client, err := s3.NewMinioClient(cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Endpoint)
	if err != nil {
		return fmt.Errorf("could not create s3 client: %w", err)
	}

	var mailClient mail.Sender
	if cfg.Env == config.EnvDevelopment {
		mailClient = mail.NewStdoutSender()
	} else {
		mailClient = mail.NewSendgridSender(cfg.SendgridApiKey)
	}
	h := handler.ContestHandler{
		ContestantRepo: contestantrepo,
		EntryRepo:      entryrepo,
		S3:             s3Client,
		MailClient:     mailClient,
		Config:         *cfg,
	}

	r := gin.Default()
	r.Static("/assets", "./frontend")
	r.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	r.GET("/", h.HandleGet)
	r.POST("/", h.HandlePost)
	r.GET("/confirm/:token", h.HandleGetConfirm)
	r.GET("/success/:contestantid", h.HandlePostSuccess)
	r.NoRoute(h.HandleNotFound)
	err = r.Run(":8080")
	return err
}
