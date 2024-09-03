# generate swagger
make swagger:
	swag init

# run server
make run:
	go run main.go