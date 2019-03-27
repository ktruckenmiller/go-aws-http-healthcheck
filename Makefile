exec_name=aws_app_health

build-linux:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(exec_name) main.go

build-mac:
	go build -o $(exec_name) main.go

develop: build-develop
	docker run -it --rm \
	-v $(shell pwd):$(shell pwd) \
	-w $(shell pwd) \
	-e GOOS=linux \
	-e GOARCH=386 \
	-e CGO_ENABLED=0 \
	-e AWS_REGION=us-west-2 \
	-e IAM_ROLE \
	aws_app_health sh

build-develop:
	docker build -t aws_app_health --target BUILDER .

build-docker:
	docker build -t aws_app_health
