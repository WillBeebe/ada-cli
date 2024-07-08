install:
	go build
	mv ada $(HOME)/.local/bin

build:
	GOOS=linux GOARCH=amd64 go build -o ada-linux-amd64
