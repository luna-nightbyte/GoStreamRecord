

# Get the latest Git tag for versioning
GIT_TAG := $(shell git describe --tags --always --dirty)


# DEV

.PHONY: run
run gen-loutput-key:
	mkdir -p output/internal/app
	cp -r internal/app/* output/internal/app
	go build \
		-ldflags="-X 'GoStreamRecord/internal/db.Version=$(GIT_TAG)'" \
		-o ./output/server main.go && \
	cd output && \
	./server


.PHONY: gen-output-key
gen-loutput-key:
	rm -rf /output/.env
	echo "SESSION_KEY=$(shell head -c 32 /dev/urandom | base64)" >> ./output/.env


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

