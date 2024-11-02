clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./internal/...

security:
	gosec ./...

test: clean critic security
	go test -v -timeout 180s -coverprofile=cover.out -cover ./internal/... ./test/...
	go tool cover -func=cover.out