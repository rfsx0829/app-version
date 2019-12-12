elf: main.go
	CGO_ENABLED=0 GOOS=linux go build -o main -ldflags '-extldflags "-static"' main.go

run: 
	go run main.go

build:
	go build -o main main.go

dbuild: main dockerfile
	docker build -t cloud .

drun:
	docker run -d -p 8999:8999 cloud:latest
