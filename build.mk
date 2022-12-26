.PHONY: build clean

NAME := gosak
ENV := local
GOOS := linux
GOARCH := amd64
HASH := $$(git rev-parse --short --verify HEAD)
DATE := $$(date -u '+%Y%m%dT%H%M%S')
GOVERSION = $$(go version)

build: $(NAME).$(ENV).$(GOOS).$(GOARCH)

$(NAME).$(ENV).$(GOOS).$(GOARCH):
	GOOS=$(GOOS) GOARCH=$(GOARCH) \
	    go build -tags=$(ENV) \
	        -o $(NAME) \
	        -ldflags "-X \"github.com/tienvu461/gosak/utils.Version=$(DATE)-$(HASH)\" \
	-X \"github.com/tienvu461/gosak/utils.Hash=$(HASH)\" \
	            -X \"github.com/tienvu461/gosak/utils.Date=$(DATE)\" \
	            -X \"github.com/tienvu461/gosak/utils.GoVersion=$(GOVERSION)\"" \
	        .
	mv $(NAME) $@

clean:
	-rm -rf $(NAME).$(ENV).$(GOOS).$(GOARCH)
