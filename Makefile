default: install

deps:
	@go get github.com/tools/godep

bindir:
	@mkdir -p bin/

build: bindir deps
	GOOS=linux GOARCH=amd64 go build -o bin/ec2-env-linux-amd64 ec2-env.go

install: build
	@cp bin/* ${GOPATH}/bin

release: build
	@aws s3 cp bin/ec2-env-linux-amd64 s3://opsee-releases/go/ec2-env/ec2-env-linux-amd64

clean:
	@rm -f bin/*
