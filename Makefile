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