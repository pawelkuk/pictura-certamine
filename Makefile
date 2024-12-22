.PHONY: migrate_contest_up migrate_contest_down run
migrate_contest_up:
	migrate -path pkg/domain/contest/db/migration/ \
	 -database "sqlite3://data/pictura-certamine.db?x-migrations-table=contest_migrations" \
	 -verbose up

migrate_contest_down:
	migrate -path pkg/domain/contest/db/migration/ \
	 -database "sqlite3://data/pictura-certamine.db?x-migrations-table=contest_migrations" \
	 -verbose down -all
	rm ./**/pictura-certamine.db

run: bundle templ
	go run cmd/main.go

templ:
	templ generate

debug: templ
	dlv debug cmd/main.go

bundle:
	cd frontend && npm run build && cd ..