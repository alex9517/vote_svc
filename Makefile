NAME		:=	Votes Servcie
VERSION		:=	$(shell cat ./VERSION)
BINARY		:=	vote-svc
SOURCE_DIR	:=	.
BUILD_DATE	:=	$(shell date '+%Y-%m-%d %H:%M:%S')
INSTALL_DIR	:=	/home/alex/bin

IMAGE_STD	:=	ms-vote-svc-std
IMAGE_SLIM	:=	ms-vote-svc-slim
IMAGE_ALPINE	:=	ms-vote-svc-alpine

DOCKERFILE_STD		:=	Dockerfile
DOCKERFILE_SLIM		:=	Dockerfile.slim
DOCKERFILE_ALPINE	:=	Dockerfile.alpine


export GODEBUG=x509ignoreCN=0

.PHONY: build run test test_integ test_coverage install clean lint

.DEFAULT_GOAL: build

# build: $(SOURCE)

build:
	@echo "Building project.."
	@echo $(VERSION) / $(BUILD_DATE)
	go build -o $(BINARY) cmd/main.go

image:
#	It creates image using 'debian:bullseye' and an existing binary.
#	That is, you must first create binary 'barcode-create', and then
#	you can run this procedure.
	docker build -t $(IMAGE_STD) -f $(DOCKERFILE_STD) .

image_slim:
#	It builds binary, then creates image based on 'debian:bullseye-slim'.
	docker build -t $(IMAGE_SLIM) -f $(DOCKERFILE_SLIM) .

image_alpine:
#	It builds binary, then creates image based on 'alpine:latest'.
	docker build -t $(IMAGE_ALPINE) -f $(DOCKERFILE_ALPINE) .

server:
	go run cmd/main.go

test:
	./runtests.sh

#test_integ:
#	go test -tags=integration

test_coverage:
	go test ./... -coverprofile=coverage.out

install:
	cp $(BINARY) $(INSTALL_DIR)

clean:
	go clean
	rm -f $(BINARY)

dep:
	go mod download

fmt:
	go fmt ./...
	goimports -w $(SOURCE)

vet:
	go vet ./...

lint:
	golint $(SOURCE)
