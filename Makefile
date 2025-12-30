.PHONY: test coverage coverage-html

test:
	cd backend && go test ./...

coverage:
	cd backend && go test -coverprofile=coverage.out ./...
	cd backend && go tool cover -func=coverage.out

coverage-html:
	cd backend && go tool cover -html=coverage.out
