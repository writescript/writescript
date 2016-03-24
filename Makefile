# writescript Makefile

all: test build

build:
	@cd cmd/writescript && go build
	@./cmd/writescript/writescript -v

test:
	@golint
	@go test -bench=. -benchmem -v

install:
	@cd cmd/writescript && go install
