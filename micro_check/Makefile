VERSION=test
REGISTRY=hub.test.cn
NAME=micro_check

build:
	GOARCH=amd64 GOOS=linux go build -o $(NAME)

build-image:
	GOARCH=amd64 GOOS=linux go build -o $(NAME)
	docker build -f Dockerfile -t $(REGISTRY)/micro/$(NAME):$(VERSION) .
	docker push $(REGISTRY)/micro/$(NAME):$(VERSION)
	rm -rf $(NAME)
