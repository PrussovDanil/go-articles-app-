DB_URL=postgresql://gouser:gopass@localhost:5432/go_article_app?sslmode=disable

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

.PHONY: migrate-reset
migrate-reset:
	migrate -path migrations -database "$(DB_URL)" down
	migrate -path migrations -database "$(DB_URL)" up

.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

.PHONY: run
run:
	go run main.go

.PHONY: test-db
test-db:
	PGPASSWORD=gopass psql -h localhost -U gouser -d go_article_app -c "SELECT COUNT(*) FROM users;"
	PGPASSWORD=gopass psql -h localhost -U gouser -d go_article_app -c "SELECT COUNT(*) FROM articles;"