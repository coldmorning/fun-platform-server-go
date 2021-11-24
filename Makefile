test: 
	go test -v -cover -covermode=atomic ./...
run-env:
        docker rm -f $(docker ps -a -q)
	docker image rm fun-platform-server
	docker-compose up --build -d
run-main:
	go run main.go
