#####################################################################################
## print usage information
help:
	@echo 'Usage:'
	@cat ${MAKEFILE_LIST} | grep -e "^## " -A 1 | grep -v '\-\-' | sed 's/^##//' | cut -f1 -d":" | \
		awk '{info=$$0; getline; print "  " $$0 ": " info;}' | column -t -s ':' | sort 
.PHONY: help
#####################################################################################
## call units tests
test/unit:
	go test -race -count 1 ./...	
.PHONY: test/unit
#####################################################################################
## code vet and lint
test/lint: 
	go vet `go list ./... | grep -v mocks`
	go install golang.org/x/lint/golint@latest
	golint -set_exit_status ./...
.PHONY: test/lint
#####################################################################################
## cleans prepared data for dockeriimage generation
clean:
	go mod tidy
	go clean