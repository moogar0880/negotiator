# run the ci directives. used primarily by travis when vetting a PR
.PHONY: ci
ci: vendor test

# run unit tests and generate a coverage report
.PHONY: coverage
coverage:
	go test -v -coverprofile=coverage.out

# run unit tests and generate a coverage report, then open the HTML version of
# the coverage report in the browser
.PHONY: coverage.view
coverage.view: coverage
	go tool cover -html=coverage.out

# generate a local copy of the package's godoc and opern them in the browser
.PHONY: godoc
godoc:
	open "http://localhost:6060/pkg/github.com/moogar0880/negotiator/" && \
	godoc -http=:6060

# run unit tests sans coverage reports
.PHONY: test
test:
	go test -v -cover

# pull in vendored dependencies into the vendor/ directory
.PHONY: vendor
vendor:
	glide install
