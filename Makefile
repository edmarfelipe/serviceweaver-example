generate:
	weaver generate ./...

build: generate
	go build -o bin/sw-blog

run-single: build
	go run .

run-multi: build
	weaver multi deploy weaver.toml