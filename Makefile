.PHONY: migrate_contest_up migrate_contest_down
migrate_contest_up:
	migrate -path pkg/domain/contest/db/migration/ \
	 -database "sqlite3://data/pictura-certamine.db?x-migrations-table=contest_migrations" \
	 -verbose up

migrate_contest_down:
	migrate -path pkg/domain/contest/db/migration/ \
	 -database "sqlite3://data/pictura-certamine.db?x-migrations-table=contest_migrations" \
	 -verbose down