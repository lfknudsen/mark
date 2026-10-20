all:
	go build
	install

gen:
	go generate ./...