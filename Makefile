VERSION="0.0"
GIT_COMMIT=$(shell git rev-list -1 HEAD)

SOURCES=$(shell ls cmd/gadget/*.go)

DEPENDS=\
	golang.org/x/crypto/ssh\
	github.com/tmc/scp\
	gopkg.in/yaml.v2\
	github.com/satori/go.uuid\
	golang.org/x/crypto/ssh\
	golang.org/x/crypto/ssh/terminal\

gadget: $(SOURCES)
	go build -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./cmd/gadget

tidy:
	go fmt ./cmd/gadget
test: gadget
	mkdir test-project
	./gadget -C test-project init
	./gadget -C test-project build
	./gadget -C test-project deploy
	./gadget -C test-project start
	./gadget -C test-project status
	./gadget -C test-project logs
	./gadget -C test-project stop
	./gadget -C test-project status
	./gadget -C test-project delete

get:
	go get ${DEPENDS}
