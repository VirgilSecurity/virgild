.PHONY: get test test_all test_integration build clear build build_artifacts
docker: get test_all build_docker clear

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

build_docker:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build  --ldflags '-extldflags "-static"' -o build/linux_x64/$(ARTF)
	docker build -t $(IMAGENAME) --build-arg GIT_COMMIT=$(GIT_COMMIT) --build-arg GIT_BRANCH=$(GIT_BRANCH) .

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
