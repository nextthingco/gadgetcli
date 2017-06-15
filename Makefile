VERSION="0.0"
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date --iso-8601)
VERSION_FILE=libgadget/version.go

GADGET_SOURCES=$(shell ls gadgetcli/*.go)
GADGETOSINIT_SOURCES=$(shell ls gadgetosinit/*.go)
LIBGADGET_SOURCES=$(shell ls libgadget/*.go)

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
	github.com/gin-gonic/gin\

gadget: genversion $(GADGET_SOURCES) $(VERSION_FILE) $(LIBGADGET_SOURCES)
	@echo "Building Gadget"
	@rm -rf ${GOPATH}/src/github.com/nextthingco/libgadget
	@cp -r libgadget ${GOPATH}/src/github.com/nextthingco/
	@go install -ldflags="-X libgadget.Version=$(VERSION) -X libgadget.GitCommit=$(GIT_COMMIT)" -v github.com/nextthingco/libgadget
	@go build -o gadget -ldflags="-s -w" -v ./gadgetcli

genversion:
	@echo "package libgadget" > $(VERSION_FILE)
	@echo "const (" >> $(VERSION_FILE)
	@echo "	Version = \"${VERSION}\"" >> $(VERSION_FILE)
	@echo "	GitCommit = \"${GIT_COMMIT}\"" >> $(VERSION_FILE)
	@echo "	BuildDate = \"${BUILD_DATE}\"" >> $(VERSION_FILE)
	@echo ")" >> $(VERSION_FILE)

gadget_release: $(GADGET_SOURCES) $(LIBGADGET_SOURCES)
	@echo "Building Gadget Release"
	@mkdir -p build/linux
	@mkdir -p build/linux_arm
	@mkdir -p build/linux_arm64
	@mkdir -p build/windows
	@mkdir -p build/darwin
	@GOOS=linux GOARCH=amd64 go build -o build/linux/gadget -ldflags="-s -w" -v ./gadgetcli
	@GOOS=linux GOARCH=arm go build -o build/linux_arm/gadget -ldflags="-s -w" -v ./gadgetcli
	@GOOS=linux GOARCH=arm64 go build -o build/linux_arm64/gadget -ldflags="-s -w" -v ./gadgetcli
	@GOOS=windows GOARCH=amd64 go build -o build/windows/gadget.exe -ldflags="-s -w" -v ./gadgetcli
	@GOOS=darwin GOARCH=amd64 go build -o build/darwin/gadget -ldflags="-s -w" -v ./gadgetcli

gadgetosinit_release: $(GADGETOSINIT_SOURCES) $(LIBGADGET_SOURCES)
	@echo "Building Gadget Release"
	@mkdir -p build/linux_arm
	@mkdir -p build/linux_arm64
	@GOOS=linux GOARCH=arm go build -o build/linux_arm/gadgetosinit -ldflags="-s -w" ./gadgetosinit
	@GOOS=linux GOARCH=arm go build -o build/linux_arm64/gadgetosinit -ldflags="-s -w" ./gadgetosinit

tidy:
	@echo "Tidying up sources"
	@go fmt ./gadgetcli
	@go fmt ./gadgetosinit
	@go fmt ./libgadget

test: $(GADGET_SOURCES) $(GADGET_SOURCES)
	@echo "Testing Gadget"
	@rm -f /tmp/gadget.yml gadgetcli/gadget.yml
	@go test -ldflags="-s -w" -v ./gadgetcli
	@go test -ldflags="-s -w" -v ./libgadget

get:
	@echo "Downloading external dependencies"
	@go get ${DEPENDS}
