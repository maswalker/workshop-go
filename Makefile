.PHONY: all
all: clean api

.PHONY: help
help:
	@echo "all                   Clean and build all"
	@echo "clean                 Clean"
	@echo "api                   Build api"

.PHONY: clean
clean:
	rm -rf build

.PHONY: api
api:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-server main.go
	upx -9 build/api-server
