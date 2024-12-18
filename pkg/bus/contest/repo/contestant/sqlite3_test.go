package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/mail"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestSqliteRepoContestantCRUD(t *testing.T) {
	skipCI(t)
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "./data/pictura-certamine.db")
	if err != nil {
		panic(err)
	}
	repo := SQLiteRepo{DB: db}
	cc := model.Contestant{
		ID:             "abcd",
		Email:          mail.Address{Address: "test@example.com"},
		FirstName:      "Alice",
		Surname:        "Doe",
		Birthdate:      time.Now(),
		PolicyAccepted: true,
	}
	err = repo.Create(ctx, &cc)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("contest created:", cc.ID)
	cc.PolicyAccepted = false
	err = repo.Update(ctx, &cc)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	contestRead := model.Contestant{ID: cc.ID}
	err = repo.Read(ctx, &contestRead)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if contestRead.PolicyAccepted {
		log.Fatalf("policy_accepted not updated")
	}
	err = repo.Delete(ctx, &cc)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
