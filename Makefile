

# Get the latest Git tag for versioning
GIT_TAG := $(shell git describe --tags --always --dirty)
FULL_TAG := $(GIT_TAG)-$(shell git describe --all)

user := $1
pass := $2



# GOLANG
#-- app
.PHONY: app
app:
	mkdir -p output/internal/app
	cp -r internal/app/* output/internal/app
	go build \
		-ldflags="-X 'GoStreamRecord/internal/db.Version=$(FULL_TAG)'" \
		-o ./output/server main.go 
	cd output && \
	./server
	
#-- reset password
.PHONY: reset-pwd
reset-pwd:
	go build \
		-ldflags="-X 'GoStreamRecord/internal/db.Version=$(FULL_TAG)'" \
		-o ./server main.go 
	./server reset-pwd $(user) $(pass)
	
	
# DOCKER
.PHONY: build
build: build-base push-base build-app push-app

.PHONY: base
build-base: 
	docker build \
		--build-arg TAG=$(GIT_TAG) \
		-t lunanightbyte/gorecord-base:$(GIT_TAG) . \
		-f ./docker/Dockerfile.base \
	

.PHONY: build-app
build-app: push-app
	docker build \
		--build-arg TAG=$(GIT_TAG) \
		-t lunanightbyte/gorecord:$(GIT_TAG) . \
		-f ./docker/Dockerfile.run

.PHONY: push-app
push-app:
	docker push lunanightbyte/gorecord:$(GIT_TAG)

.PHONY: push-app
push-base:
	docker push lunanightbyte/gorecord-base:$(GIT_TAG)
