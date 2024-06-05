image:
	docker build -t pingatus .

build:
	go build -o pingatus -mod=vendor