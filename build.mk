.PHONY: build build-all clean

NAME := gosak
ENV := local
GOOS := linux
GOARCH := amd64
TAG := $$(git tag --points-at HEAD)
HASH := $$(git rev-parse --short --verify HEAD)
DATE := $$(date -u '+%Y%m%dT%H%M%S')
GOVERSION = $$(go version)
ALLOS = linux darwin
ALLARCH = amd64 arm64

build: $(NAME).$(GOOS).$(GOARCH)

$(NAME).$(GOOS).$(GOARCH):
	GOOS=$(GOOS) GOARCH=$(GOARCH) \
	    go build \
	        -o $(NAME) \
	        -ldflags "-X \"github.com/tienvu461/gosak/utils.Version=$(TAG)-$(HASH)-$(DATE)\" \
	            -X \"github.com/tienvu461/gosak/utils.Hash=$(HASH)\" \
	            -X \"github.com/tienvu461/gosak/utils.Date=$(DATE)\" \
				-X \"github.com/tienvu461/gosak/utils.Arch=$(GOARCH)\" \
				-X \"github.com/tienvu461/gosak/utils.Os=$(GOOS)\" \
	            -X \"github.com/tienvu461/gosak/utils.GoVersion=$(GOVERSION)\"" \
	        .
	mv $(NAME) $@

clean:
	-rm -rf $(NAME).*

build-all:
	for os in $(ALLOS); do \
	    for arch in $(ALLARCH); do \
	        GOOS=$$os GOARCH=$$arch \
	        go build \
	            -o $(NAME).$$os.$$arch \
	            -ldflags "-X \"github.com/tienvu461/gosak/utils.Version=$(TAG)-$(HASH)-$(DATE)\" \
	                -X \"github.com/tienvu461/gosak/utils.Hash=$(HASH)\" \
	                -X \"github.com/tienvu461/gosak/utils.Date=$(DATE)\" \
	                -X \"github.com/tienvu461/gosak/utils.Arch=$$arch\" \
	                -X \"github.com/tienvu461/gosak/utils.Os=$$os\" \
	                -X \"github.com/tienvu461/gosak/utils.GoVersion=$(GOVERSION)\"" \
	                .;\
	    done \
	done
