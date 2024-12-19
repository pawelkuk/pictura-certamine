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
	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/model"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestSqliteRepoContestCRUD(t *testing.T) {
	skipCI(t)
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "./data/pictura-certamine.db")
	if err != nil {
		panic(err)
	}
	repo := SQLiteRepo{DB: db}
	cc := model.Contest{
		ID:       "abcd",
		Name:     "best picture",
		Slug:     model.ParseSlug("best picture"),
		Start:    time.Now(),
		End:      time.Now(),
		IsActive: false,
	}
	err = repo.Create(ctx, &cc)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("contest created:", cc.ID)
	cc.End = cc.End.Add(24 * time.Hour)
	err = repo.Update(ctx, &cc)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	contestRead := model.Contest{ID: cc.ID}
	err = repo.Read(ctx, &contestRead)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = repo.Delete(ctx, &cc)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
