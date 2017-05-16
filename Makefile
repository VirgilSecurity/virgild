.PHONY: get test test_all clean_artifact build_artifacts docker build_docker docker_test docker_inspect docker_dockerhub_publish

PROJECT =virgild
IMAGENAME=$(PROJECT)
DOCKERHUB_REPOSITORY=virgilsecurity/$(IMAGENAME)

ifeq ($(OS),Windows_NT)
TARGET_OS ?= windows
else
TARGET_OS ?= $(shell uname -s | tr A-Z a-z)
endif

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

all: get build
.DEFAULT_GOAL := all


define tag_docker
  @if [ "$(GIT_BRANCH)" = "master" ]; then \
    docker tag $(IMAGENAME) $(1):latest; \
  fi
  @if [ "$(GIT_BRANCH)" != "master" ]; then \
    docker tag $(IMAGENAME) $(1):$(GIT_BRANCH); \
  fi
endef

clean:
	rm $(BUILD_FILE_NAME)

clean_artifact:
		rm -rf artf

test: get
		go test -v ./...

test_all: get
	go test -v ./... -tags=integration

test_coverage:
	@echo "" > coverage.txt

	@for d in $$(go list ./... | grep -v vendor); do \
	    go test -race -tags=integration -coverprofile=profile.out -covermode=atomic $$d; \
	    if [ -f profile.out ]; then \
	        cat profile.out >> coverage.out; \
	        rm profile.out; \
	    fi; \
	done


$(GOPATH)/src/gopkg.in/virgilsecurity/virgil-crypto-go.v4/virgil_crypto_go.go:
ifeq ($(C_CRYPTO),true)
	go get -d gopkg.in/virgilsecurity/virgil-crypto-go.v4
	cd $$GOPATH/src/gopkg.in/virgilsecurity/virgil-crypto-go.v4 ;	 make
endif

get:$(GOPATH)/src/gopkg.in/virgilsecurity/virgil-crypto-go.v4/virgil_crypto_go.go
	go get github.com/jteeuwen/go-bindata/...
	go get -v -d -t -tags docker  ./...

build:
	go-bindata -pkg db -o modules/card/db/bindata.go -prefix modules/card/ modules/card/db/migrations/
	CGO_ENABLED=1 GOOS=$(TARGET_OS) go build  $(BUILD_ARGS) -o $(BUILD_FILE_NAME)

build_in_docker-env:
ifeq ($(TARGET_OS),linux)
	make
else
	docker pull virgilsecurity/virgil-crypto-go-env
	docker run -it --rm -v "$$PWD":/go/src/github.com/VirgilSecurity/virgild -w /go/src/github.com/VirgilSecurity/virgild virgilsecurity/virgil-crypto-go-env make
endif

docker: build_docker docker_test

build_docker: build_in_docker-env
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
	go test -tags=docker -v
	docker-compose down

docker_dockerhub_publish:
	$(call tag_docker, $(DOCKERHUB_REPOSITORY))
	docker push $(DOCKERHUB_REPOSITORY)

docker_inspect:
		docker inspect -f '{{index .ContainerConfig.Labels "git-commit"}}' $(IMAGENAME)
		docker inspect -f '{{index .ContainerConfig.Labels "git-branch"}}' $(IMAGENAME)

build_artifacts: clean_artifact build
	mkdir -p artf/src/$(PROJECT)
	mv $(BUILD_FILE_NAME) artf/src/$(PROJECT)/

ifeq ($(TARGET_OS),windows)
	cd artf/src &&	zip -r ../$(ARTF_OS_NAME)-amd64.zip . &&	cd ../..
else
	tar -zcvf artf/$(ARTF_OS_NAME)-amd64.tar.gz -C artf/src .
endif

	rm -rf artf/src
