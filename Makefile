USER:=$(shell id -u)
GROUP:=$(shell id -g)

deploy.app:
	cd frontend; npm install && npx mix build -p
	make build
	ansible-playbook -i deploy/hosts.yml deploy/server.yml -t configuration -e @deploy/vars/server.yml -e "USER=1000" -e "GROUP=1000" --ask-vault-pass
init:
	ansible-playbook -i deploy/hosts.yml deploy/local.yml -t configuration -e @deploy/vars/local.yml -e "USER=$(USER)" -e "GROUP=$(GROUP)"
build:
	docker compose run --rm app sh -c "CGO_ENABLED=1 go build -ldflags '-linkmode external -w -extldflags \"-static\"' -o tmp/app cmd/main.go"
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