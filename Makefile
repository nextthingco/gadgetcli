SOURCES=\
	cmd/gadget/main.go \
	cmd/gadget/config.go \
	cmd/gadget/build.go \
	cmd/gadget/init.go \
	cmd/gadget/start.go

DEPENDS=\
	golang.org/x/crypto/ssh\
	github.com/tmc/scp\
	gopkg.in/yaml.v2\
	github.com/satori/go.uuid\
	golang.org/x/crypto/ssh\
	golang.org/x/crypto/ssh/terminal\

gadget: $(SOURCES)
	go build -ldflags="-s -w" -v ./cmd/gadget

test: gadget
	./gadget -C tests build

get:
	go get ${DEPENDS}
