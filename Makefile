# Get the latest Git tag for versioning
TAG := $(shell git describe --tags --always --dirty)

.PHONY: build
build: base app

.PHONY: base
base:
	docker build \
		--build-arg TAG=$(TAG) \
		-t lunanightbyte/gorecord-base:$(TAG) .
	docker push lunanightbyte/gorecord-base:$(TAG) # Must be pushed t use in 'app' image

.PHONY: app
app:
	docker build \
		--build-arg TAG=$(TAG) \
		-t lunanightbyte/gorecord:$(TAG) .

.PHONY: push
push:
	docker push lunanightbyte/gorecord:$(TAG)

# Clean up dangling images
.PHONY: clean
clean:
	docker image prune -f


