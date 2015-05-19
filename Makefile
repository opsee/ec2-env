default: install

deps:
	@go get github.com/tools/godep

bindir:
	@mkdir -p bin/

build: bindir deps
	@go build -o bin/ec2-env ec2-env.go

install: build
	@cp bin/* ${GOPATH}/bin

clean:
	@rm -f bin/*
