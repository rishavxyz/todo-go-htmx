dev:
	$(GOPATH)/bin/air

build:
	go build -o bin/todo-app ./main.go

run: build
	./bin/todo-app
