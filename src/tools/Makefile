-include Makefile.options
#####################################################################################
## print usage information
help:
	@echo 'Usage:'
	@cat ${MAKEFILE_LIST} | grep -e "^## " -A 1 | grep -v '\-\-' | sed 's/^##//' | cut -f1 -d":" | \
		awk '{info=$$0; getline; print "  " $$0 ": " info;}' | column -t -s ':' | sort 
.PHONY: help
#####################################################################################
## generate mock objects for test
generate/mocks: 
	go install github.com/petergtz/pegomock/v3/...@v3.0.1
	go generate ./...
.PHONY: generate/mocks
#####################################################################################
## call units tests
test/unit:
	go test -race -count 1 ./...	
.PHONY: test/unit
#####################################################################################
## call snyk tests
test/snyk:
	@docker run --rm -i$(tty) --env SNYK_TOKEN=$(SNYK_TOKEN) -v $$(pwd):/app snyk/snyk:golang
.PHONY: test/snyk
#####################################################################################
../../test-reports:
	mkdir -p $@
## generates test reports
test/report: | ../../test-reports
	go install github.com/jstemmer/go-junit-report@latest
	go test ./... -v -race 2>&1 | go-junit-report > ../../test-reports/report.xml
.PHONY: test/report
#####################################################################################
## code vet and lint
test/lint: 
	go vet ./...
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run -v ./...
.PHONY: test/lint
#####################################################################################
## cleans prepared data for dockeriimage generation
clean:
	go mod tidy
	go clean