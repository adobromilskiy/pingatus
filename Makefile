image:
	docker build -t pingatus .

build:
	go build -o pingatus -mod=vendor

test:
	go test -cover -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html