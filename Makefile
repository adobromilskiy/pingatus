image:
	docker build -t pingatus .

build:
	go build -o pingatus -mod=vendor

test:
	docker run -tid --name testmongo -p 27117:27017 mongo:7.0
	-go test -cover -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	docker rm -f testmongo