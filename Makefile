SOURCES=\
	cmd/gadget/main.go \
	cmd/gadget/config.go \
	cmd/gadget/build.go \
	cmd/gadget/start.go

gadget: $(SOURCES)
	go build -v ./cmd/gadget

test: gadget
	./gadget -C tests build