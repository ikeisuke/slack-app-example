.PHONY: all test clean build run

all: install test build prepare package deploy

install:
	go mod download

build:
	GOARCH=amd64 GOOS=linux go build -o build/slash/slash handlers/slash/main.go
	GOARCH=amd64 GOOS=linux go build -o build/event/event handlers/event/main.go

test:
	go test handlers/slash/main.go
	go test handlers/event/main.go

clean:
	-rm -rf build

prepare:
	aws cloudformation deploy \
	 	--stack-name slack-command-example-prepare \
	 	--template-file cfn/prepare.yaml \
	 	--no-fail-on-empty-changeset

package:
	aws cloudformation package \
	 	--s3-bucket `aws cloudformation describe-stacks \
                    		--stack-name slack-command-example-prepare \
                    		--query "Stacks[0].Outputs[?OutputKey=='PackageBucketName'].OutputValue | [0]" \
                    		--output text` \
	 	--template-file cfn/application.yaml \
	 	--output-template-file build/application.yaml

deploy:
	aws cloudformation deploy \
	 	--stack-name slack-command-example-application \
	 	--template-file build/application.yaml \
	 	--capabilities CAPABILITY_IAM \
	 	--no-fail-on-empty-changeset \
	 	--parameter-overrides \
	 		SlackSigningSecret=${SLACK_SIGNING_SECTET} \
	 		EnvironmentEncryptionKeyArn=${ENVIRONMENT_ENCRYPTION_KEY_ARN}
