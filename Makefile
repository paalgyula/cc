# The remote address is for testing only. You should set it via
remote ?= "127.0.0.1:5000"
LDFLAGS := "-s -w -X 'github.com/paalgyula/cc/config.Remote=$(remote)'"

default: clean cc

help:           ## Show this help.
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v grep -F | sed -e 's/\\$$//' | sed -e 's/##//'

clean:		## Removes all generated binaries
	rm -f cc
	rm -f cc_tiny

tinybuild:  ## Builds this project with tinygo's docker container
	# docker run --rm -v $(PWD):/src -w /src tinygo/tinygo:0.28.1 tinygo build -o cc_tiny -tags "netgo tiny" -gc=leaking -no-debug ./
	docker run --rm -v $(PWD):/src -w /src tinygo/tinygo:0.28.1 tinygo build -o cc_tiny -tags tiny ./

cc: ## Makes the cc binary for linux
	set -e
	GOOS=linux ARCH=x64 go build -trimpath -tags netgo,nodebug -ldflags=$(LDFLAGS) -o cc .
	upx -9 --no-backup cc
	ls -lah cc
