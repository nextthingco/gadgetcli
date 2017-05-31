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
	github.com/sirupsen/logrus\

gadget: $(SOURCES)
	@echo "Building Gadget"
	@go build -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./cmd/gadget

tidy:
	@echo "Tidying up sources"
	@go fmt ./cmd/gadget

test: $(SOURCES)
	@echo "Testing Gadget"
	@go test -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./cmd/gadget

get:
	@echo "Downloading external dependencies"
	@go get ${DEPENDS}
