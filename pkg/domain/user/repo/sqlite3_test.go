package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	model "github.com/pawelkuk/pictura-certamine/pkg/domain/user/model"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestSqliteRepoCRUD(t *testing.T) {
	skipCI(t)
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "./data/todo_test.db")
	if err != nil {
		panic(err)
	}
	repo := SQLiteRepo{DB: db}

	uu, err := model.Parse("bob@hacker.com", model.WithPassword("password"))
	if err != nil {
		panic(err)
	}
	err = repo.Create(ctx, uu)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("user created:", uu.Email)
	err = repo.Update(ctx, uu)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	userRead := model.User{ID: uu.ID}
	err = repo.Read(ctx, &userRead)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = repo.Delete(ctx, uu)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func TestSqliteRepoQuery(t *testing.T) {
	skipCI(t)
	ctx := context.Background()
	err := os.Remove("./data/todo.db")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
	}
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	db, err := sql.Open("sqlite3", "./data/todo.db")
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:///%s", "/Users/pawkuk/workspace/todo/pkg/db/migration"),
		"sqlite3", driver,
	)

	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	repo := SQLiteRepo{DB: db}

	u1, err := model.Parse("bob@hacker.com", model.WithPassword("password"))
	if err != nil {
		panic(err)
	}
	err = repo.Create(ctx, u1)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	u2, err := model.Parse("bob2@hacker.com", model.WithPassword("password2"))
	if err != nil {
		panic(err)
	}
	err = repo.Create(ctx, u2)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	usersAll, err := repo.Query(ctx, model.QueryFilter{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(usersAll) != 2 {
		log.Fatalf("len(usersAll) != 2: %d", len(usersAll))
	}
	usersBob, err := repo.Query(ctx, model.QueryFilter{Email: &mail.Address{Address: `bob%`}})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(usersBob) != 2 {
		log.Fatalf("len(usersBob) != 2: %d", len(usersBob))
	}
	users, err := repo.Query(ctx, model.QueryFilter{ID: &u1.ID})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(users) != 1 {
		log.Fatalf("len(users) != 1: %d", len(users))
	}
}
