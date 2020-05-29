GOPATH:=$(shell pwd)
GO:=go
GOFLAGS:=-v -p 1

default:  autoendpoint remover test injector

all: clean default

autoendpoint: bin/autoendpoint

test: bin/test_sanitize_header

remover: bin/sub_remover

injector: bin/sub_injector


bin/autoendpoint: src/empowerthings.com/autoendpoint/autoendpoint.go
	@echo "========== Compiling $@ =========="
	sh -c 'export GOPATH=${GOPATH}; $(GO) install $(GOFLAGS) empowerthings.com/autoendpoint/ '

bin/test_sanitize_header: src/empowerthings.com/autoendpoint/test/test_sanitize_header/test_sanitize_header.go
	@echo "========== Compiling $@ =========="
	sh -c 'export GOPATH=${GOPATH}; $(GO) install $(GOFLAGS) empowerthings.com/autoendpoint/test/test_sanitize_header'

bin/sub_remover: src/empowerthings.com/sub_remover/sub_remover.go
	@echo "========== Compiling $@ =========="
	sh -c 'export GOPATH=${GOPATH}; $(GO) install $(GOFLAGS) empowerthings.com/sub_remover/'


bin/sub_injector: src/empowerthings.com/sub_injector/sub_injector.go
	@echo "========== Compiling $@ =========="
	sh -c 'export GOPATH=${GOPATH}; $(GO) install $(GOFLAGS) empowerthings.com/sub_injector/'

deploy:
	tar -cvjf "rep-autopush-`git describe --tags --abbrev=0`.tar.bz2" bin

VERSION_DOCKER=$(shell echo `git describe`.`date +"%Y%m%d%H%M%S"` > .version)
docker: ${VERSION_DOCKER}
	docker build . -t rep-autopush-$$(git describe --abbrev=0)

clean:
	@echo "Deleting generated binary files ..."; sh -c 'if [ -d bin ]; then  find bin/ -type f -exec rm {} \; -print ; fi; rm -Rf bin'
	@echo "Deleting generated archive files ..."; sh -c 'if [ -d pkg ]; then  find pkg -type f -name \*.a -exec rm {} \; -print ; fi;  rm -Rf pkg'
	@echo "Deleting emacs backup files ..."; find . -type f -name \*~ -exec rm {} \; -print
	@echo "Deleting log files ..."; find . -maxdepth 1 -type f \( -name \*.log.\* -o -name \*.log \) -exec rm {} \; -print
