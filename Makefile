

# Get the latest Git tag for versioning
GIT_TAG := $(shell git describe --tags)
FULL_TAG := $(GIT_TAG)-$(shell git describe --all)
BINARY_NAME=$(GIT_TAG)

user := $1
pass := $2
LENGHT=32


# GOLANG
#-- app
.PHONY: app
app:
	go build \
		-ldflags="-X 'GoStreamRecord/internal/db.Version=$(FULL_TAG)'" \
		-ldflags="-X 'GoStreamRecord/internal/cli.BinaryName=$(BINARY_NAME)'" \
		-o ./output/server main.go 
	cd output && \
	$(BINARY_NAME)
	
#-- reset password
.PHONY: reset-pwd
reset-pwd:
	go build \
		-ldflags="-X 'GoStreamRecord/internal/db.Version=$(FULL_TAG)'" \
		-ldflags="-X 'GoStreamRecord/internal/cli.BinaryName=$(BINARY_NAME)'" \
		-o $(BINARY_NAME) main.go 
	$(BINARY_NAME) reset-pwd $(user) $(pass)
	
	
# DOCKER
.PHONY: build
build: build-base push-base build-app push-app

.PHONY: base
base: 
	docker build \
		--build-arg TAG=$(GIT_TAG) \
		-t lunanightbyte/gorecord-base:$(GIT_TAG) . \
		-f ./docker/Dockerfile.base \
	

.PHONY: build-app
build-app: 
	docker build \
		--build-arg VERSION=$(GIT_TAG) \
		--build-arg BINARY=$(BINARY_NAME) \
		-t lunanightbyte/gorecord:$(GIT_TAG) . \
		-f ./docker/Dockerfile.run

.PHONY: push-app
push-app:
	docker push lunanightbyte/gorecord:$(GIT_TAG)

.PHONY: push-base
push-base:
	docker push lunanightbyte/gorecord-base:$(GIT_TAG)



.PHONY: new-cookie-token
new-cookie-token:
	go build \
		-ldflags="-X 'GoStreamRecord/internal/cli.BinaryName=$(BINARY_NAME)'" \
		-ldflags="-X 'GoStreamRecord/internal/db.Version=$(FULL_TAG)'" \
		-o ./output/server main.go 
	cd output && \
	$(BINARY_NAME) gen-cookie-token $(LENGHT)
	
	

.PHONY: new-session-token
new-session-token:
	go build \
		-ldflags="-X 'GoStreamRecord/internal/cli.BinaryName=$(BINARY_NAME)'" \
		-ldflags="-X 'GoStreamRecord/internal/db.Version=$(FULL_TAG)'" \
		-o ./output/$(BINARY_NAME) main.go 
	cd output && \
	./$(BINARY_NAME) gen-session-token $(LENGHT)
	
