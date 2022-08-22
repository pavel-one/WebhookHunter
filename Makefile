include ${PWD}/.env
export

up:
	docker-compose up -d && make log
down:
	docker-compose down
exec:
	docker-compose exec app bash
migrate:
	docker-compose exec app migrate create -ext sql -dir db/migrations ${name}
migrate.up:
	docker-compose exec app migrate -database "postgres://$(DB_PASSWORD):$(DB_USER)@$(DB_HOST):5432/dev?sslmode=disable" -path db/migrations up
migrate.down:
	docker-compose exec app migrate -database "postgres://$(DB_PASSWORD):$(DB_USER)@$(DB_HOST):5432/dev?sslmode=disable" -path db/migrations down
exec.root:
	docker-compose exec -u root app bash
log:
	docker-compose logs -f app