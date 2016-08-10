.PHONY: coverage coverage.view test

coverage:
	go test -v -coverprofile=coverage.out

coverage.view: coverage
	go tool cover -html=coverage.out

test:
	go test -v -cover
