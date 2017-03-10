.PHONY: get test test_all clear_artifact build_artifacts docker_test docker_dockerhub_publish

PROJECT =virgild
IMAGENAME=$(PROJECT)
DOCKERHUB_REPOSITORY=virgilsecurity/$(IMAGENAME)
TARGET_OS ?= $(shell uname -s | tr A-Z a-z)

ifeq ($(TARGET_OS),darwin)
ARTF_OS_NAME?=macosx
else
ARTF_OS_NAME?=$(TARGET_OS)
endif

ifeq ($(TARGET_OS),windows)
BUILD_FILE_NAME?=$(PROJECT).exe
C_CRYPTO=false
else
BUILD_FILE_NAME?=$(PROJECT)
C_CRYPTO?=true
endif

BUILD_ARGS=
ifeq ($(C_CRYPTO),true)
BUILD_ARGS+=-tags=c_crypto
endif
ifneq ($(TARGET_OS),darwin)
BUILD_ARGS+= --ldflags '-extldflags "-static"'
endif

.DEFAULT_GOAL := $(BUILD_FILE_NAME)


define tag_docker
  @if [ "$(GIT_BRANCH)" = "master" ]; then \
    docker tag $(IMAGENAME) $(1):latest; \
  fi
  @if [ "$(GIT_BRANCH)" != "master" ]; then \
    docker tag $(IMAGENAME) $(1):$(GIT_BRANCH); \
  fi
endef

clear_artifact:
		rm -rf artf

test: get
		go test -v ./...

test_all: get
	go test -v ./... -tags=integration


$(GOPATH)/src/gopkg.in/virgilsecurity/virgil-crypto-go.v4/virgil_crypto_go.go:
ifeq ($(C_CRYPTO),true)
	go get -d gopkg.in/virgilsecurity/virgil-crypto-go.v4
	cd $$GOPATH/src/gopkg.in/virgilsecurity/virgil-crypto-go.v4 ;	 make
endif

get: $(GOPATH)/src/gopkg.in/virgilsecurity/virgil-crypto-go.v4/virgil_crypto_go.go
	go get -v -d -t -tags docker  ./...


$(BUILD_FILE_NAME): get
	CGO_ENABLED=1 GOOS=$(TARGET_OS) go build  $(BUILD_ARGS) -o $(BUILD_FILE_NAME)


docker: build_docker docker_test

build_docker: build
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

docker_dockerhub_publish:
	$(call tag_docker, $(DOCKERHUB_REPOSITORY))
	docker push $(DOCKERHUB_REPOSITORY)

docker_inspect:
		docker inspect -f '{{index .ContainerConfig.Labels "git-commit"}}' $(IMAGENAME)
		docker inspect -f '{{index .ContainerConfig.Labels "git-branch"}}' $(IMAGENAME)

build_artifacts: clear_artifact $(BUILD_FILE_NAME)
	mkdir -p artf/src/$(PROJECT)
	mv $(BUILD_FILE_NAME) artf/src/$(PROJECT)/

ifeq ($(TARGET_OS),windows)
	cd artf/src &&	zip -r ../$(ARTF_OS_NAME)-amd64.zip . &&	cd ../..
else
	tar -zcvf artf/$(ARTF_OS_NAME)-amd64.tar.gz -C artf/src .
endif

	rm -rf artf/src
