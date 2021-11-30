test: 
	go test -v -cover -covermode=atomic ./...
docker-clean:
        docker rm -f $(docker ps -a -q)
	docker image rm fun-platform-server
run-env:
	docker-compose up --build -d
run-main:
	go run main.go
