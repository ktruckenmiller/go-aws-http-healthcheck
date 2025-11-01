exec_name=bootstrap

# Get version SHA from git, fallback to timestamp if not in git repo
VERSION_SHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo $(shell date +%s))

build-linux:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(exec_name) ./...

build-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(exec_name) ./...
	chmod +x $(exec_name)
	zip bootstrap-arm64-$(VERSION_SHA).zip $(exec_name)
	zip bootstrap-arm64.zip $(exec_name)
	rm $(exec_name)

build-mac: build-develop
	docker run -it --rm \
	--platform linux/arm64 \
	-v ${PWD}:${PWD} \
	-w ${PWD} \
	$(exec_name) \
	go mod download && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build  -ldflags="-w -s" -o $(exec_name) -tags lambda.norpc main.go
	chmod +x $(exec_name)
	zip bootstrap-arm64-$(VERSION_SHA).zip $(exec_name)
	zip bootstrap-arm64.zip $(exec_name)
	rm $(exec_name)

build-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(exec_name) -tags lambda.norpc main.go
	chmod +x $(exec_name)
	zip bootstrap-amd64-$(VERSION_SHA).zip $(exec_name)
	zip bootstrap-amd64.zip $(exec_name)
	rm $(exec_name)

upload:
	@echo "Uploading with version SHA: $(VERSION_SHA)"
	aws s3 cp bootstrap-amd64-$(VERSION_SHA).zip s3://kloudcover-public-us-west-2-601394826940/healthcheck/bootstrap-amd64-$(VERSION_SHA).zip
	aws s3 cp bootstrap-arm64-$(VERSION_SHA).zip s3://kloudcover-public-us-west-2-601394826940/healthcheck/bootstrap-arm64-$(VERSION_SHA).zip
	aws s3 cp bootstrap-amd64.zip s3://kloudcover-public-us-west-2-601394826940/healthcheck/bootstrap-amd64.zip
	aws s3 cp bootstrap-arm64.zip s3://kloudcover-public-us-west-2-601394826940/healthcheck/bootstrap-arm64.zip

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
	@echo "Deploying with version SHA: $(VERSION_SHA)"
	aws cloudformation deploy \
		--role-arn arn:aws:iam::601394826940:role/cdk-hnb659fds-cfn-exec-role-601394826940-us-west-2 \
		--template-file stack.yml \
		--stack-name go-http-healthcheck \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			ServiceName=my-ip \
			ServiceURL=https://my-ip.clustermaestro.com \
			Environment=prod \
			PhoneNumber=${PHONE_NUMBER} \
			VersionSHA=$(VERSION_SHA)
