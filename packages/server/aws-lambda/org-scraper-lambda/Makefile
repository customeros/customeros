CURRENT_DIR := $(notdir $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST))))))
TAG := $(CURRENT_DIR)

print_info:
	@echo "CURRENT_DIR:" $(CURRENT_DIR)
	@echo "TAG:" $(TAG)

build:
	docker build --platform=linux/amd64 -t $(TAG) .

tag:
	docker tag $(TAG):latest 769325097132.dkr.ecr.eu-west-1.amazonaws.com/$(TAG)

push:
	docker push 769325097132.dkr.ecr.eu-west-1.amazonaws.com/$(TAG):latest

update:
	aws lambda update-function-code --function-name $(TAG) --image-uri 769325097132.dkr.ecr.eu-west-1.amazonaws.com/$(TAG):latest --query 'LastModified'

ecr-login:
	aws ecr get-login-password --region eu-west-1 | docker login --username AWS --password-stdin 769325097132.dkr.ecr.eu-west-1.amazonaws.com

ecr-create:
	aws ecr create-repository --repository-name $(TAG)

iam-role-create:
	aws iam create-role --role-name $(TAG)-role --assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}'
	# aws iam attach-role-policy --role-name $(TAG)-role --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole

aws-lambda-create:
	aws lambda create-function \
  		--function-name $(TAG) \
  		--package-type Image \
  		--code ImageUri=769325097132.dkr.ecr.eu-west-1.amazonaws.com/$(TAG):latest \
  		--role arn:aws:iam::769325097132:role/$(TAG)-role \
  		#--vpc-config SubnetIds=subnet-09cf54e0a71976027,subnet-0f06f3a9119904d59,subnet-0e7267afbefb930da,subnet-0e559edf3be5811ec,subnet-00e7a15aa1dd530a9,subnet-0f9c88bb7aed55cc3,SecurityGroupIds=sg-0c766973f8bfd7397,sg-04c3ee495c3ca66a7 \
  		--timeout 5 \
  		--memory-size 512


all: build tag push update