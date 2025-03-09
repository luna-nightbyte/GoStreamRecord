

# Get the latest Git tag for versioning
GIT_TAG := $(shell git describe --tags --always --dirty)


# DEV

.PHONY: build-go
build-go:
	go build \
		-ldflags="-X 'GoStreamRecord/internal/db.Version=$(GIT_TAG)'" \
		-o ./output/server main.go 

.PHONY: run
run:
	mkdir -p output/internal/app
	cp -r internal/app/* output/internal/app
	make build-go
	cd output && \
	./server

# DOCKER
.PHONY: build
build: base app

.PHONY: base
base: push-base
	docker build \
		--build-arg TAG=$(GIT_TAG) \
		-t lunanightbyte/gorecord-base:$(GIT_TAG) . \
		-f ./docker/Dockerfile.base \

.PHONY: app
app: push-app
	docker build \
		--build-arg TAG=$(GIT_TAG) \
		-t lunanightbyte/gorecord:$(GIT_TAG) . \
		-f ./docker/Dockerfile.run \
	docker push lunanightbyte/gorecord:$(GIT_TAG)

.PHONY: push-app
push-app:
	docker push lunanightbyte/gorecord:$(GIT_TAG)

.PHONY: push-app
push-base:
	docker push lunanightbyte/gorecord-base:$(GIT_TAG)

# Clean up dangling images
.PHONY: clean
clean:
	docker image prune -f

