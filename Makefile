.PHONY: test coverage coverage-html

test:
	docker compose exec api go test ./...

coverage:
	docker compose exec api go test -coverprofile=coverage.out ./...
	docker compose exec api go tool cover -func=coverage.out

coverage-html: coverage
	cd backend && go tool cover -html=coverage.out
