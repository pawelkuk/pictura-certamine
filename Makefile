.PHONY: migrate_contest_up migrate_contest_down
migrate_contest_up:
	migrate -path pkg/bus/contest/db/migration/ \
	 -database "sqlite3://data/todo.db?x-migrations-table=contest_migrations" \
	 -verbose up

migrate_contest_down:
	migrate -path pkg/bus/contest/db/migration/ \
	 -database "sqlite3://data/todo.db?x-migrations-table=contest_migrations" \
	 -verbose down