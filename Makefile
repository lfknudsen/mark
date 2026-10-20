all:
	go build -C cmd/mark
	cd cmd/mark; install

gen:
	go generate ./...