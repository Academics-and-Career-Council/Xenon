run:
	go run main.go

build:
	go build -o build/rogue main.go 

install:
    go mod download
    go build -o build/rogue main.go 