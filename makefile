dev:
	go run ./cmd/rinha.go

build: clean deps
	CGO_ENABLED=0 go build -o ./bin/nell_challenge ./cmd/main.go

deps:
	go mod tidy

clean:
	rm -rf ./bin
.
docker-build:
	docker build -t caiotony/nell_challenge-go 

docker-start:
	docker-compose up --build -d

docker-down:
	docker-compose down -v --remove-orphans
