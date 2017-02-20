.PHONY: get test test_all test_integration build clear build build_artifacts

ARTF = virgild

get:
	go get -v -t  ./...

test_all: test test_integration

test:
	go test -v ./...

test_integration:
	go test -v ./... -tags=integration

clear:
	rm -rf build

build:
	go build -o $(ARTF)

build_artifacts: clear get
	CGO_ENABLED=1 GOARCH=amd64 go build  --ldflags '-extldflags "-static"' -o build/linux_x64/$(ARTF)
#	CGO_ENABLED=1 GOARCH=386 go build  --ldflags '-extldflags "-static"' -o build/linux_x86/$(ARTF)
#	CGO_ENABLED=1 GOARCH=arm go build  --ldflags '-extldflags "-static"' -o build/linux_arm/$(ARTF)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build  --ldflags '-extldflags "-static"' -o build/windows_x86/$(ARTF).exe
#	CGO_ENABLED=1 GOOS=windows GOARCH=386 CC=x86_64-w64-mingw32-gcc go build  --ldflags '-extldflags "-static"' -o build/windows_x64/$(ARTF).exe
	GOOS=linux go build -v --ldflags '-extldflags "-static"' -o build/linux/$(ARTF)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC=/osxcross/target/bin/o64-clang go build --ldflags '-extldflags "-static"' -o build/macos/$(ARTF)
