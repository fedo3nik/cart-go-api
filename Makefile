.PHONY: go_lint postgres_run migrate_up migrate_down swagger_spec swagger_ui

go_lint:
	docker run --rm -v ${PWD}:/app -w /app/ golangci/golangci-lint:v1.36-alpine golangci-lint run -v --timeout=5m

postgres_run:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=spoon -e PGDATA=/var/lib/postgresql/data/pgdata -v /custom/mount:/var/lib/postgresql/data -d postgres:9.6.20

migrate_up:
	migrate -source file://internal/infrastructure/database/migrations/ -database "postgresql://postgres:password@localhost:5432/cart_api?sslmode=disable" up

migrate_down:
	migrate -source file://internal/infrastructure/database/migrations/ -database "postgresql://postgres:password@localhost:5432/cart_api?sslmode=disable" down

swagger_spec:
	swagger generate spec -o ./swagger.yml

swagger_ui:
	swagger serve -F=swagger swagger.yml