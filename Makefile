GO=go
PRODUCTION_OS=linux
PRODUCTION_ARCH=386 #x86
PRODUCTION_PORT=$(shell cat conf/library.toml |grep port|cut -d " " -f 3)
TAG=nil

help: ## This help
	$(info Available Targets)
	$(info -----------------)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

doc: ## Cleans and creates the swagger documentation
	$(MAKE) clean
	$(MAKE) swagger 

image: ## Creates a docker image for application deployment 
	$(MAKE) doc 
	env GOOS=$(PRODUCTION_OS) GOARCH=$(PRODUCTION_ARCH) $(GO) build
	docker build --build-arg http_port=$(PRODUCTION_PORT) -t $(TAG) \
		-f build/Dockerfile .

compose-up: ## Runs/restarts the application locally using docker containers 
	docker-compose -f build/deploy/docker-compose.yml up -d --force-recreate

compose-dn: ## Stops the application
	docker-compose -f build/deploy/docker-compose.yml down
	
app: ## Cleans and builds swagger documentation and the application 
	$(MAKE) doc
	go build

swagger: ## Create swagger documentation
	swag init --md docs -o docs/swagger

clean: ## Clean all generated binaries and code
	go clean 
	rm -rf docs/swagger/*

.PHONY: doc image compose-up compose-down dev help clean
