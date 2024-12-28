package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pawelkuk/pictura-certamine/pkg/domain/auth/model"
	"github.com/robfig/cron/v3"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestSqliteRepoPeriodicTaskCRUD(t *testing.T) {
	skipCI(t)
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "./data/todo_test.db")
	if err != nil {
		panic(err)
	}
	repo := SQLiteRepo{DB: db}
	schedule := "* * * * *"
	_, err = cron.ParseStandard(schedule)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	tt := model.Session{
		Token:  model.SessionToken{Value: "xyz"},
		UserID: 1,
		Expiry: time.Now(),
	}
	err = repo.Create(ctx, &tt)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("periodictask created:", tt.Token.Value)
	tt.Expiry = tt.Expiry.Add(24 * time.Hour)
	err = repo.Update(ctx, &tt)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	taskRead := model.Session{Token: tt.Token}
	err = repo.Read(ctx, &taskRead)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = repo.Delete(ctx, &tt)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
