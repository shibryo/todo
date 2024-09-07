# setup
setup:
	go mod vendor

# generate swagger
swagger:
	swag init -g cmd/todoApp/main.go

# run server
run:
	make swagger
	go run cmd/todoApp/main.go

# run test
test:
	go test -v ./...

# golangci-lint
lint:
	golangci-lint run

# coverage html
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

