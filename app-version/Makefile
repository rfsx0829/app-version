build: main.go
	CGO_ENABLED=0 GOOS=linux go build -o main -ldflags '-extldflags "-static"' main.go

run: 
	go run main.go

dbuild: main dockerfile
	docker build -t app-version .

drun:
	docker run -d -p 8000:8000 app-version:latest
