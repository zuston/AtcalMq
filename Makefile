GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

PULL_BINARY_NAME=./target/pull

PUSH_BINARY_NAME=./target/push

CONSOLE_BINARY_NAME=./target/console

BACKUP_BINARY_NAME=./target/backup


all: pull push console
pull:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(PULL_BINARY_NAME) -v main.go

push:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(PUSH_BINARY_NAME) -v push.go

console:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(CONSOLE_BINARY_NAME) -v newconsole.go

backup:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BACKUP_BINARY_NAME) -v backup.go
