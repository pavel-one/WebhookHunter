USER:=$(shell id -u)
GROUP:=$(shell id -g)

init:
	ansible-playbook -i deploy/hosts.yml deploy/local.yml -t configuration -e @deploy/vars/local.yml -e "USER=$(USER)" -e "GROUP=$(GROUP)"
build:
	docker compose run --rm app sh -c "CGO_ENABLED=1 go build -ldflags '-linkmode external -w -extldflags \"-static\"' -o tmp/app cmd/main.go"
build.test:
	docker compose run --rm app sh -c "CGO_ENABLED=1 go build -ldflags '-linkmode external -w -extldflags \"-static\"' -o tmp/test cmd/main.go"
run:
	docker compose run --rm app sh -c "CGO_ENABLED=1 go build -ldflags '-linkmode external -w -extldflags \"-static\"' -o tmp/app cmd/main.go && ./tmp/app"
up:
	docker-compose up -d && make log
down:
	docker-compose down
exec:
	docker-compose exec app sh
exec.root:
	docker-compose exec -u root app sh
log:
	docker-compose logs -f app