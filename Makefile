exec_name=bootstrap

build-linux:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(exec_name) ./...

build-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(exec_name) ./...
	chmod +x $(exec_name)
	zip my-lambda-function.zip $(exec_name)

build-mac: build-develop
	docker run -it --rm \
	--platform linux/arm64 \
	-v ${PWD}:${PWD} \
	-w ${PWD} \
	$(exec_name) \
	go mod download && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build  -ldflags="-w -s" -o $(exec_name) -tags lambda.norpc main.go
	chmod +x $(exec_name)
	zip build.zip $(exec_name)
	rm $(exec_name)

develop: build-develop
	docker run -it --rm \
	-v $(shell pwd):$(shell pwd) \
	-w $(shell pwd) \
	-e GOOS=linux \
	-e GOARCH=386 \
	-e CGO_ENABLED=0 \
	-e AWS_REGION=us-west-2 \
	-e IAM_ROLE \
	$(exec_name) sh
	#URL=https://my-ip.clustermaestro.com REGION=us-east-1 METRIC_NAME=boston go run main.go

build-develop:
	docker build -t $(exec_name) --target builder --platform linux/arm64 .

build-docker:
	docker build -t $(exec_name) .

deploy:
	aws cloudformation deploy \
		--role-arn arn:aws:iam::601394826940:role/cdk-hnb659fds-cfn-exec-role-601394826940-us-west-2 \
		--template-file stack.yml \
		--stack-name go-http-healthcheck \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			ServiceName=my-ip \
			ServiceURL=https://my-ip.clustermaestro.com \
			Environment=prod \
			PhoneNumber=${PHONE_NUMBER}
