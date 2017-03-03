.PHONY: get test test_all test_integration build clear build build_artifacts build_docker
docker: get test_all build_docker docker_test

ARTF =virgild
IMAGENAME=$(ARTF)
DOCKERHUB_REPOSITORY=virgilsecurity/$(IMAGENAME)
OS = $(shell uname -s)

define tag_docker
  @if [ "$(GIT_BRANCH)" = "master" ]; then \
    docker tag $(IMAGENAME) $(1):latest; \
  fi
  @if [ "$(GIT_BRANCH)" != "master" ]; then \
    docker tag $(IMAGENAME) $(1):$(GIT_BRANCH); \
  fi
endef

get:
	go get -v -d -t -tags docker  ./...
ifeq ($(strip $(OS)),Linux)
	wget https://cdn.virgilsecurity.com/crypto-go/virgil-crypto-2.0.4-go-linux-x86_64.tgz -P $$GOPATH/src/gopkg.in/virgilsecurity/virgil-crypto-go.v4/
	tar -xvf $$GOPATH/src/gopkg.in/virgilsecurity/virgil-crypto-go.v4/virgil-crypto-2.0.4-go-linux-x86_64.tgz --strip-components=1 -C $$GOPATH/src/gopkg.in/virgilsecurity/virgil-crypto-go.v4/
endif

test_all: test test_integration

test:
	go test -v ./...

test_integration:
	go test -v ./... -tags=integration

clear:
	rm -rf build

build:
	go build -o $(ARTF)

build_in_docker:
	go get -v  ./...
	CGO_ENABLED=1 GOARCH=amd64 go build -tags=c_crypto  --ldflags '-extldflags "-static"' -o build/docker/$(ARTF)

build/docker/$(ARTF):
	docker build -t build_docker -f build_docker .
	docker run --rm -v "$$PWD":/go/src/github.com/VirgilSecurity/virgild -w /go/src/github.com/VirgilSecurity/virgild build_docker make build_in_docker

build_docker: build/docker/$(ARTF)
	docker build -t $(IMAGENAME) --build-arg GIT_COMMIT=$(GIT_COMMIT) --build-arg GIT_BRANCH=$(GIT_BRANCH) .

docker_test:
ifndef SYNC_TOKEN
	$(error SYNC_TOKEN is undefined. Access token for get access to Virgil cloud)
endif
ifndef SYNC_APP_ID
	$(error SYNC_APP_ID is undefined. Application card id in Virgil cloud)
endif
ifndef SYNC_APP_KEY
	$(error SYNC_APP_KEY is undefined. Private key in base64 incoding)
endif
ifndef SYNC_APP_KEY_PASS
	$(error SYNC_APP_KEY_PASS is undefined. Password for private key)
endif

	# CACHE
	docker-compose up -d virgild_cache
	go test -tags=docker -run Cache -v
	docker-compose down
	# SYNC
	docker-compose up -d virgild_sync
	go test -tags=docker -run Sync -v
	docker-compose down
	# LOCAL
	docker-compose up -d virgild_local
	go test -tags=docker -run Local -v
	docker-compose down

docker_dockerhub_tag:
	$(call tag_docker, $(DOCKERHUB_REPOSITORY))

docker_dockerhub_push:
	docker push $(DOCKERHUB_REPOSITORY)

docker_inspect:
		docker inspect -f '{{index .ContainerConfig.Labels "git-commit"}}' $(IMAGENAME)
		docker inspect -f '{{index .ContainerConfig.Labels "git-branch"}}' $(IMAGENAME)

build_artifacts: get build/inux-amd64.tar.gz build/windows-amd64.zip build/macosx-amd64.tar.gz

build/inux-amd64.tar.gz: build/docker/$(ARTF)
	mkdir -p build/linux-amd64/$(ARTF)
	mv build/docker/$(ARTF) build/linux-amd64/$(ARTF)/$(ARTF)
	tar -zcvf build/linux-amd64.tar.gz -C build/linux-amd64/ .
#	CGO_ENABLED=1 GOARCH=386 go build  --ldflags '-extldflags "-static"' -o build/linux_x86/$(ARTF)
#	CGO_ENABLED=1 GOARCH=arm go build  --ldflags '-extldflags "-static"' -o build/linux_arm/$(ARTF)

build/windows-amd64.zip:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build  --ldflags '-extldflags "-static"' -o build/windows-amd64/$(ARTF)/$(ARTF).exe

	cd build/windows-amd64 &&	zip -r ../windows-amd64.zip . &&	cd ../..
#	CGO_ENABLED=1 GOOS=windows GOARCH=386 CC=x86_64-w64-mingw32-gcc go build  --ldflags '-extldflags "-static"' -o build/windows_x64/$(ARTF).exe

build/macosx-amd64.tar.gz:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC=o64-clang go build -tags=c_crypto  -o build/macosx-amd64/$(ARTF)/$(ARTF)
	tar -zcvf build/macosx-amd64.tar.gz -C build/macosx-amd64/ .
