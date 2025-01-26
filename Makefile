.PHONY: migrate_contest_up migrate_contest_down run migrate_auth_up migrate_auth_down migrate_user_up migrate_user_down migrate_up migrate_down
migrate_contest_up:
	migrate -path pkg/domain/contest/db/migration/ \
	 -database "sqlite3://data/pictura-certamine.db?x-migrations-table=contest_migrations" \
	 -verbose up

migrate_contest_down:
	migrate -path pkg/domain/contest/db/migration/ \
	 -database "sqlite3://data/pictura-certamine.db?x-migrations-table=contest_migrations" \
	 -verbose down -all

migrate_auth_up:
	migrate -path pkg/domain/auth/db/migration/ \
	 -database "sqlite3://data/pictura-certamine.db?x-migrations-table=auth_migrations" \
	 -verbose up

migrate_auth_down:
	migrate -path pkg/domain/auth/db/migration/ \
	 -database "sqlite3://data/pictura-certamine.db?x-migrations-table=auth_migrations" \
	 -verbose down -all

migrate_user_up:
	migrate -path pkg/domain/user/db/migration/ \
	 -database "sqlite3://data/pictura-certamine.db?x-migrations-table=user_migrations" \
	 -verbose up

migrate_user_down:
	migrate -path pkg/domain/user/db/migration/ \
	 -database "sqlite3://data/pictura-certamine.db?x-migrations-table=user_migrations" \
	 -verbose down -all

migrate_up: migrate_auth_up migrate_user_up migrate_contest_up
	echo "migrate_up done" 

migrate_down: migrate_auth_down migrate_user_down migrate_contest_down
	rm ./**/pictura-certamine.db
	echo "migrate_down done"

run: templ
	go run cmd/app/main.go

templ: bundle
	templ generate

debug: templ
	dlv debug cmd/app/main.go

bundle:
	cd frontend && npm run build && cd ..

build-image:
	docker build -t pawelkuk/pictura-certamine:0.0.0 .

push-image:
	docker push pawelkuk/pictura-certamine:0.0.0 

rsync:
	rsync -r ./frontend/dist ubuntu@c5:/home/ubuntu/workspace/frontend/dist &&\
		rsync -r ./.prod.envrc ubuntu@c5:/home/ubuntu/workspace/.prod.envrc &&\
		rsync -r ./Makefile ubuntu@c5:/home/ubuntu/workspace/Makefile

deploy:
	docker pull pawelkuk/pictura-certamine:0.0.0 && source .prod.envrc && docker compose down && docker compose up -d