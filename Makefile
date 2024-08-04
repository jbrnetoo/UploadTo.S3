build: 
	@go build -o bin/upload ./main.go

test: 
	@go test -v ./...

run: build
	@./bin/upload