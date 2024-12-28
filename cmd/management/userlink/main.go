package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/mail"
	"os"

	"github.com/caarlos0/env/v11"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/config"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/user/model"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/user/repo"
)

func main() {
	err := exec()
	if err != nil {
		log.Fatal(err)
	}
}

func exec() error {
	cfg := &config.Config{}
	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("could not parse config: %w", err)
	}
	db, err := sql.Open("sqlite3", "./data/pictura-certamine.db")
	if err != nil {
		return fmt.Errorf("could not open db: %w", err)
	}
	repo := repo.SQLiteRepo{DB: db}

	if len(os.Args) != 2 {
		return fmt.Errorf("usage: <email>")
	}
	email := os.Args[1]
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("could not parse email: %w", err)
	}
	u, err := model.Parse(addr.Address, model.WithAuthorizationToken(), model.WithActivationToken())
	if err != nil {
		return fmt.Errorf("could not parse user: %w", err)
	}
	err = repo.Create(context.Background(), u)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}
	fmt.Println("user initialized...")
	link := fmt.Sprintf("%s/user/%s", cfg.BaseURL, u.AuthorizationToken)
	fmt.Printf("activation link: %s\n", link)
	return nil
}
