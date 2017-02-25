# writescript Makefile

VERSION=0.3.1

all: test release

build:
	@echo "build v${VERSION}"
	@cd cmd/writescript && go build
	@./cmd/writescript/writescript --version

test: test-src test-cli

test-src:
	@golint
	@go test -v ./...

bench:
	@go test -bench=. -benchmem -v ./...

test-cli:
	@echo "run writescript..."
	@./cmd/writescript/writescript -p docs/tutorials/1-minute/main.wjs

install:
	@cd cmd/writescript && go install

clean:
	rm -rf release
	rm -rf cmd/writescript/build

release: release-darwin release-linux release-windows release-freebsd release-netbsd
	@echo "release build finished!"

define r
	@echo "release build v${VERSION} $1 $2 $3"
	@mkdir -p build/$1/$2
	@cd cmd/writescript && env GOOS=$1 GOARCH=$2 go build -o build/$1/$2/writescript$3
	@cd cmd/writescript/build/$1/$2 && zip -r writescript_v0.3.1_$1_$2.zip writescript$3
	@mv cmd/writescript/build/$1/$2/writescript_v0.3.1_$1_$2.zip release/writescript_v0.3.1_$1_$2.zip
endef

release-darwin:
	@mkdir -p release
	$(call r,darwin,amd64)
	$(call r,darwin,386)
	$(call r,darwin,arm)
	@echo "release build darwin finished!"

release-linux:
	@mkdir -p release
	$(call r,linux,amd64)
	$(call r,linux,386)
	$(call r,linux,arm)
	@echo "release build linux finished!"

release-windows:
	@mkdir -p release
	$(call r,windows,amd64,.exe)
	$(call r,windows,386,.exe)
	@echo "release build windows finished!"

release-freebsd:
	@mkdir -p release
	$(call r,freebsd,amd64)
	$(call r,freebsd,386)
	$(call r,freebsd,arm)
	@echo "release build freebsd finished!"

release-netbsd:
	@mkdir -p release
	$(call r,netbsd,amd64)
	$(call r,netbsd,386)
	$(call r,netbsd,arm)
	@echo "release build netbsd finished!"

.PHONY: build
