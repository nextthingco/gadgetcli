VERSION="0.0"
GIT_COMMIT=$(shell git rev-list -1 HEAD)

SOURCES=$(shell ls gadgetcli/*.go)

DEPENDS=\
	golang.org/x/crypto/ssh\
	github.com/tmc/scp\
	gopkg.in/yaml.v2\
	github.com/satori/go.uuid\
	golang.org/x/crypto/ssh\
	golang.org/x/crypto/ssh/terminal\
	github.com/sirupsen/logrus\
	gopkg.in/cheggaaa/pb.v1\
	github.com/nextthingco/logrus-gadget-formatter\

gadget: $(SOURCES)
	@echo "Building Gadget"
	@go build -o gadget -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./gadgetcli

release: $(SOURCES)
	@echo "Building Gadget Release"
	@mkdir -p build/linux
	@mkdir -p build/linux_arm
	@mkdir -p build/linux_arm64
	@mkdir -p build/windows
	@mkdir -p build/darwin
	@GOOS=linux GOARCH=amd64 go build -o build/linux/gadget -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./gadgetcli
	@GOOS=linux GOARCH=arm go build -o build/linux_arm/gadget -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./gadgetcli
	@GOOS=linux GOARCH=arm64 go build -o build/linux_arm64/gadget -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./gadgetcli
	@GOOS=windows GOARCH=amd64 go build -o build/windows/gadget.exe -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./gadgetcli
	@GOOS=darwin GOARCH=amd64 go build -o build/darwin/gadget -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./gadgetcli

tidy:
	@echo "Tidying up sources"
	@go fmt ./gadgetcli

test: $(SOURCES)
	@echo "Testing Gadget"
	@go test -ldflags="-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -v ./gadgetcli

get:
	@echo "Downloading external dependencies"
	@go get ${DEPENDS}
