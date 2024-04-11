.PHONY: install
install:
	go install -v ./.
	make tidy

.PHONY: build
build:
	go build -o bin/go-coco
	make tidy

.PHONY: tidy
tidy:
	go mod tidy