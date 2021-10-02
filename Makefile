test: 
	go test -v -cover -covermode=atomic ./...
run-env:
	docker-compose up --build -d
run:
	go run main.go