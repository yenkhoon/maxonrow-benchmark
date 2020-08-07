export GO111MODULE = on

all: deps build install

deps:
	go mod vendor

build:
	go build $(LDFLAGS) $(TAGS) -mod vendor -o ./build/benchmark ./benchmark.go
	
install:
	go install $(LDFLAGS) $(TAGS) -mod vendor ./benchmark.go
	
