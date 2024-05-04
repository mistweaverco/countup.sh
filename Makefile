BIN_NAME=countup

build-windows-64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X 'main.VERSION=$(VERSION)'" -o dist/windows/$(BIN_NAME).exe cmd/main.go
build-linux-64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X 'main.VERSION=$(VERSION)'" -o dist/linux/$(BIN_NAME) cmd/main.go
build-macos-arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-s -w -X 'main.VERSION=$(VERSION)'" -o dist/macos/$(BIN_NAME) cmd/main.go

builds: build-linux-64 build-macos-arm64 build-windows-64

archives:
	zip --junk-paths -r dist/countup.sh-$(VERSION)-windows.zip dist/windows
	tar -czvf dist/countup.sh-$(VERSION)-linux.tar.gz -C dist/linux .
	tar -czvf dist/countup.sh-$(VERSION)-macos.tar.gz -C dist/macos .

release:
	gh release create --generate-notes v$(VERSION) dist/countup.sh-$(VERSION)-linux.tar.gz dist/countup.sh-$(VERSION)-macos.tar.gz dist/countup.sh-$(VERSION)-windows.zip

run:
	go run -ldflags "-X 'main.VERSION=development'" cmd/main.go "Dummy Timer"
