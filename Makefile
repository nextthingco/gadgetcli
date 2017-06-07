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
	gopkg.in/cheggaaa/pb.v1\

gadget: $(SOURCES)
	@echo "Building Gadget"
	@go build -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./cmd/gadget

release: $(SOURCES)
	@echo "Building Gadget Release"
	@mkdir -p build/linux
	@mkdir -p build/linux_arm
	@mkdir -p build/linux_arm64
	@mkdir -p build/windows
	@mkdir -p build/darwin
	@GOOS=linux GOARCH=amd64 go build -o build/linux/gadget -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./cmd/gadget
	@GOOS=linux GOARCH=arm go build -o build/linux_arm/gadget -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./cmd/gadget
	@GOOS=linux GOARCH=arm64 go build -o build/linux_arm64/gadget -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./cmd/gadget
	@GOOS=windows GOARCH=amd64 go build -o build/windows/gadget.exe -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./cmd/gadget
	@GOOS=darwin GOARCH=amd64 go build -o build/darwin/gadget -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./cmd/gadget

tidy:
	@echo "Tidying up sources"
	@go fmt ./cmd/gadget

test: $(SOURCES)
	@echo "Testing Gadget"
	@go test -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./cmd/gadget

get:
	@echo "Downloading external dependencies"
	@go get ${DEPENDS}
