build: main.go
	CGO_ENABLED=0 GOOS=linux go build -o main -ldflags '-extldflags "-static"' main.go

run: 
	go run main.go
