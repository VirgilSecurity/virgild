.PHONY: get test test_all test_integration build clear build build_artifacts
docker: get test_all build_docker docker_test

ARTF =virgild
IMAGENAME=$(ARTF)
DOCKERHUB_REPOSITORY=virgilsecurity/$(IMAGENAME)

define tag_docker
  @if [ "$(GIT_BRANCH)" = "master" ]; then \
    docker tag $(IMAGENAME) $(1):latest; \
  fi
  @if [ "$(GIT_BRANCH)" != "master" ]; then \
    docker tag $(IMAGENAME) $(1):$(GIT_BRANCH); \
  fi
endef

get:
	go get -v -t -tags docker  ./...

test_all: test test_integration

test:
	go test -v ./...

test_integration:
	go test -v ./... -tags=integration

clear:
	rm -rf build

build:
	go build -o $(ARTF)

build_in_docker: get
	CGO_ENABLED=1 GOARCH=amd64 go build  --ldflags '-extldflags "-static"' -o build/docker/$(ARTF)

build_docker:
	docker run --rm -v "$$PWD":/go/src/github.com/VirgilSecurity/virgild -w /go/src/github.com/VirgilSecurity/virgild golang:1.7 make build_in_docker
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
ifndef SYNC_APP_KEY
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

build_artifacts: clear get
	CGO_ENABLED=1 GOARCH=amd64 go build  --ldflags '-extldflags "-static"' -o build/linux_x64/$(ARTF)
#	CGO_ENABLED=1 GOARCH=386 go build  --ldflags '-extldflags "-static"' -o build/linux_x86/$(ARTF)
#	CGO_ENABLED=1 GOARCH=arm go build  --ldflags '-extldflags "-static"' -o build/linux_arm/$(ARTF)

	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build  --ldflags '-extldflags "-static"' -o build/windows_x86/$(ARTF).exe
#	CGO_ENABLED=1 GOOS=windows GOARCH=386 CC=x86_64-w64-mingw32-gcc go build  --ldflags '-extldflags "-static"' -o build/windows_x64/$(ARTF).exe

	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC=o64-clang go build  -o build/macos/$(ARTF)
