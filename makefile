.PHONY: vue build_go pi app run golang

# Check if programs are installed and if not set it as a dependency. 
ifneq ($(shell command -v node npm 1>/dev/null 2>&1; echo $$?), 0)
	NODE := node
endif
ifneq ($(shell command -v go 1>/dev/null 2>&1; echo $$?), 0)
	GOLANG := golang
endif
ifneq ($(shell command -v git 1>/dev/null 2>&1; echo $$?), 0)
	GIT := git
endif
ifneq ($(shell command -v curl 1>/dev/null 2>&1; echo $$?), 0)
	CURL := curl
endif
ifneq ($(shell command -v wget 1>/dev/null 2>&1; echo $$?), 0)
	WGET := wget
endif 
 
VERSION=$(shell git describe --tags)
COMMIT_HASH := $(shell git rev-parse HEAD)$(shell git diff --quiet && git diff --cached --quiet && test -z "$$(git ls-files --others --exclude-standard)" || echo "-dirty")


USERNAME=default_user
PASSWORD=default_password
ROLE=user
GROUP=viewers

pi: vue
	GOOS=linux GOARCH=arm64 go build \
	-ldflags=" \
		-X 'remoteCtrl/internal/system/version.Version=$(VERSION)' \
		-X 'remoteCtrl/internal/system/version.Shasum=$(COMMIT_HASH)'" \
	-o GoStreamRecord_Rpi

app: vue build_go

build_go: 
	go build \
	-buildvcs=false \
	-ldflags=" \
		-X 'remoteCtrl/internal/system/version.Version=$(VERSION)' \
		-X 'remoteCtrl/internal/system/version.Shasum=$(COMMIT_HASH)'" \
	-o GoStreamRecord
  
run: build_go 
	mkdir -p output/videos 
	cp ./GoStreamRecord output/GoStreamRecord
	cd output && \
	sudo ./GoStreamRecord

add-user: build_go
	mkdir -p output/settings
	mkdir -p output/videos
	cp -r --update settings/* output/settings
	cp ./GoStreamRecord output/GoStreamRecord
	cd output && \
	sudo ./GoStreamRecord add-user $(USERNAME) $(PASSWORD) $(ROLE) $(GROUP)
# Install go
.PHONY: golang
golang:
	curl -Lo /tmp/go1.21.3.linux-amd64.tar.gz \
		https://golang.org/dl/go1.21.3.linux-amd64.tar.gz
	echo "1241381b2843fae5a9707eec1f8fb2ef94d827990582c7c7c32f5bdfbfd420c8 /tmp/go1.21.3.linux-amd64.tar.gz" \
		| sha256sum --check
	sudo rm -rf /usr/local/go
	sudo tar -C /usr/local -xvzf /tmp/go1.21.3.linux-amd64.tar.gz
	sudo ln -sf /usr/local/go/bin/go /usr/local/bin/go
	sudo ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt
	sudo rm -f /tmp/go1.21.3.linux-amd64.tar.gz
	GOPATH=$(shell go env GOPATH)
 
.PHONY: vue
vue:
	cd vue/app && \
	npm install && \
	npm run build
	cd vue/login && \
	npm install && \
	npm run build