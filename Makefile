SERVICE_NAME=project-sanjiv

# go dep version https://github.com/golang/dep
DEP_VERSION=0.3.2

# check the underlying operation system and set as lowercase for further processing
OS := $(shell uname | tr A-Z a-z)

# check if dep is currently installed
DEP := $(shell command -v dep 2> /dev/null)

check_for_dep:
ifndef DEP
	@echo "go dep not installed; installing..."
	- curl -L -s https://github.com/golang/dep/releases/download/v$(DEP_VERSION)/dep-$(OS)-amd64 -o $(GOPATH)/bin/dep
	- chmod +x $(GOPATH)/bin/dep
endif
	@echo "go dep already installed."

# check for go dep and ensure that all dependencies are met
ensure_dep:
	- make check_for_dep
	- dep ensure

test:
	go test ./...

docker_run: docker_build
	docker run -d \
	    --rm -p 8080:8080 \
	    --name quote \
		-e LISTEN_PORT=8080 \
		-e LISTEN_HOST=0.0.0.0 \
		quote:latest

docker_stop:
	docker stop quote
	docker rmi quote

build_linux: ensure_dep
	GOOS=linux GOARCH=amd64 go build -o quote

docker_build: build_linux
	docker build -t quote:latest ./
