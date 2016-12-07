.PHONY: ci
ci: vendor test

.PHONY: coverage
coverage:
	go test -v -coverprofile=coverage.out

.PHONY: coverage.view
coverage.view: coverage
	go tool cover -html=coverage.out

.PHONY: test
test:
	go test -v -cover

.PHONY: vendor
vendor:
	glide install
