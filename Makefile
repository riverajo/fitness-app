.PHONY: test coverage coverage-html generate

test:
	docker compose exec api go test ./...

coverage:
	docker compose exec api go test -coverprofile=coverage.out ./...
	docker compose exec api go tool cover -func=coverage.out

coverage-html: coverage
	cd backend && go tool cover -html=coverage.out

generate:
	cd backend && go tool gqlgen generate
	cd frontend && pnpm run codegen
