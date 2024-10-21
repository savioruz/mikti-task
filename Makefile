APP_NAME=bookstore

swagger:
	swag init -g cmd/main.go

dev:
	air

build:
	LDFLAGS="-s -w" go build -o $(APP_NAME) cmd/main.go