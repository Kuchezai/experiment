include .env

run:
	go run cmd/experiment/main.go

migrate_up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose up

migrate_down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose down

linter:
	golangci-lint run --enable sqlclosecheck --enable bodyclose --enable rowserrcheck

coverage:
	go test -covermode=count -coverpkg=./... -coverprofile cover.out -v ./... && \
	go tool cover -func=cover.out

compose_up:
	docker-compose up -d --force-recreate --build

compose_down:
	docker-compose down
