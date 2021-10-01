run:
	go run main.go

build:
	go mod download
	go build -o build/xenon main.go