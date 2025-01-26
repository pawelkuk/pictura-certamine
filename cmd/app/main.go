package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v11"
	sentry "github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	authhandler "github.com/pawelkuk/pictura-certamine/pkg/domain/auth/handler"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/auth/middleware"
	auth "github.com/pawelkuk/pictura-certamine/pkg/domain/auth/repo"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/config"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/handler"
	contestant "github.com/pawelkuk/pictura-certamine/pkg/domain/contest/repo/contestant"
	entry "github.com/pawelkuk/pictura-certamine/pkg/domain/contest/repo/entry"
	crmhandler "github.com/pawelkuk/pictura-certamine/pkg/domain/crm/handler"
	contestantentry "github.com/pawelkuk/pictura-certamine/pkg/domain/crm/repo"
	userhandler "github.com/pawelkuk/pictura-certamine/pkg/domain/user/handler"
	user "github.com/pawelkuk/pictura-certamine/pkg/domain/user/repo"
	"github.com/pawelkuk/pictura-certamine/pkg/sdk/mail"
	"github.com/pawelkuk/pictura-certamine/pkg/sdk/s3"
	"golang.org/x/text/language"
)

//go:embed i18n/*
var fs embed.FS

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
	contestantentryrepo := &contestantentry.SQLiteRepo{DB: db}
	userrepo := &user.SQLiteRepo{DB: db}
	authrepo := &auth.SQLiteRepo{DB: db}

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
	contestHandler := handler.ContestHandler{
		ContestantRepo: contestantrepo,
		EntryRepo:      entryrepo,
		S3:             s3Client,
		MailClient:     mailClient,
		Config:         *cfg,
	}

	crmHandler := crmhandler.Handler{
		Repo:   contestantentryrepo,
		S3:     s3Client,
		Config: *cfg,
	}

	userHandler := userhandler.Handler{
		Repo:       userrepo,
		MailClient: mailClient,
		Config:     *cfg,
	}

	authHandler := authhandler.Handler{
		UserRepo:   userrepo,
		Repo:       authrepo,
		MailClient: mailClient,
		Config:     *cfg,
	}
	authMiddleware := middleware.Middleware{Repo: authrepo, Config: *cfg}

	r := gin.Default()

	// apply i18n middleware
	r.Use(ginI18n.Localize(ginI18n.WithBundle(&ginI18n.BundleCfg{
		DefaultLanguage:  language.Romanian,
		FormatBundleFile: "json",
		AcceptLanguage:   []language.Tag{language.Romanian, language.English, language.Polish},
		RootPath:         "./i18n/",
		UnmarshalFunc:    json.Unmarshal,
		// After commenting this line, use defaultLoader
		// it will be loaded from the file
		Loader: &ginI18n.EmbedLoader{FS: fs},
	})))

	r.GET("/i18n", func(c *gin.Context) {
		c.String(http.StatusOK, ginI18n.MustGetMessage(c, "welcome"))
	})

	r.GET("/i18n/:name", func(c *gin.Context) {
		c.String(http.StatusOK, ginI18n.MustGetMessage(
			c,
			&i18n.LocalizeConfig{
				MessageID: "welcomeWithName",
				TemplateData: map[string]string{
					"name": c.Param("name"),
				},
			}))
	})

	r.Static("/assets", "./frontend/dist")
	r.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	r.GET("/", authMiddleware.Handle, contestHandler.HandleGet)
	r.POST("/", contestHandler.HandlePost)
	r.GET("/success/:contestantid", contestHandler.HandlePostSuccess)

	r.GET("/crm", authMiddleware.Handle, crmHandler.GetAll)
	r.GET("/crm/export", authMiddleware.Handle, crmHandler.GetCSV)
	r.GET("/:env/:entryid/:filename", authMiddleware.Handle, crmHandler.GetFile)

	r.GET("/user/:authorization_token", userHandler.Get)
	r.POST("/user/:authorization_token", userHandler.Post)
	r.GET("/user/activate/:activation_token", userHandler.Activate)

	r.GET("/auth/login", authHandler.LoginGet)
	r.POST("/auth/login", authHandler.LoginPost)
	r.GET("/auth/reset", authHandler.ResetGet)
	r.POST("/auth/reset", authHandler.ResetPost)
	r.GET("/auth/password/:password_reset_token", authHandler.PasswordGet)
	r.POST("/auth/password/:password_reset_token", authHandler.PasswordPost)
	r.POST("/auth/logout/", authMiddleware.Handle, authHandler.Logout)

	r.NoRoute(contestHandler.HandleNotFound)
	err = r.Run(":8080")
	return err
}
