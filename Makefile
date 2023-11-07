BIN_DEST := /usr/bin/bm

.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p ./build/
	go build -o ./build/bm main.go

.PHONY: install
install:
	install -m 755 ./build/bm $(BIN_DEST)

.PHONY: uninstall
uninstall:
	rm -f $(BIN_DEST)