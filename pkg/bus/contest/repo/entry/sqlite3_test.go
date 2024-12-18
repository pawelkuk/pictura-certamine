package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestSqliteRepoEntryCRUD(t *testing.T) {
	skipCI(t)
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "./data/pictura-certamine.db")
	if err != nil {
		panic(err)
	}
	repo := SQLiteRepo{DB: db}
	entry := model.Entry{
		ID:           "abcdefgh",
		ContestantID: "c1",
		SessionID:    "s1",
		Status:       model.EntryStatusPending,
		ArtPieces: []model.ArtPiece{
			{Key: "/photo/1.png"}, {Key: "/photo/2.png"}, {Key: "/photo/3.png"}, {Key: "/photo/4.png"},
		},
	}
	err = repo.Create(ctx, &entry)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("entry created:", entry.ID)
	entry.Status = model.EntryStatusConfirmationEmailSent
	entry.ArtPieces = append(entry.ArtPieces, model.ArtPiece{Key: "/photo/5.png"})
	err = repo.Update(ctx, &entry)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	entryRead := model.Entry{ID: entry.ID}
	err = repo.Read(ctx, &entryRead)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(entryRead.ArtPieces) != 5 || entryRead.Status != model.EntryStatusConfirmationEmailSent {
		log.Fatalf("did not update: %v, %v", entryRead.ArtPieces, entryRead.Status)
	}
	err = repo.Delete(ctx, &entry)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
