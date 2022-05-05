.PHONY: dist dist-win dist-macos dist-linux ensure-dist-dir build install uninstall

GOBUILD=go build -ldflags="-s -w"
INSTALLPATH=/usr/local/bin

ensure-dist-dir:
	@- mkdir -p dist

dist-win: ensure-dist-dir
	# Build for Windows x64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o dist/e6dl-windows-amd64.exe main.go

dist-macos: ensure-dist-dir
	# Build for macOS x64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o dist/e6dl-darwin-amd64 main.go

	# Build for macOS ARM
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o dist/e6dl-darwin-arm64 main.go

dist-linux: ensure-dist-dir
	# Build for Linux x64
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o dist/e6dl-linux-amd64 main.go

	# Build for Linux PPC64
	GOOS=linux GOARCH=ppc64 $(GOBUILD) -o dist/e6dl-linux-ppc64 main.go

dist: dist-win dist-macos dist-linux

build:
	@- mkdir -p bin
	$(GOBUILD) -o bin/e6dl main.go
	@- chmod +x bin/e6dl

install: build
	mv bin/e6dl $(INSTALLPATH)/e6dl
	@- rm -rf bin
	@echo "e6dl was installed to $(INSTALLPATH)/e6dl. Run make uninstall to get rid of it, or just remove the binary yourself."

uninstall:
	rm $(INSTALLPATH)/e6dl

run:
	@- go run main.go