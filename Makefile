# generate swagger
swagger:
	swag init

# run server
run:
	swag init
	go run main.go

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