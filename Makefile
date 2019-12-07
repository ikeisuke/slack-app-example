.PHONY: all install test build clean prepare package deploy show-api-url sam-local-start sam-local-stop

all: install test build prepare package deploy show-api-url

install:
	go mod download

build:
	GOARCH=amd64 GOOS=linux go build -o build/slash/slash handlers/slash/main.go
	GOARCH=amd64 GOOS=linux go build -o build/event/event handlers/event/main.go
	GOARCH=amd64 GOOS=linux go build -o build/interactive/interactive handlers/interactive/main.go

test:
	go test handlers/slash/main.go
	go test handlers/event/main.go
	go test handlers/interactive/main.go

clean: sam-local-stop
	-rm -rf build

prepare:
	aws cloudformation deploy \
	 	--stack-name slack-command-example2-prepare \
	 	--template-file cfn/prepare.yaml \
	 	--no-fail-on-empty-changeset

package:
	aws cloudformation package \
	 	--s3-bucket `aws cloudformation describe-stacks \
                    		--stack-name slack-command-example2-prepare \
                    		--query "Stacks[0].Outputs[?OutputKey=='PackageBucketName'].OutputValue | [0]" \
                    		--output text` \
	 	--template-file cfn/application.yaml \
	 	--output-template-file build/application.yaml

deploy:
	aws cloudformation deploy \
	 	--stack-name slack-command-example2-application \
	 	--template-file build/application.yaml \
	 	--capabilities CAPABILITY_IAM \
	 	--no-fail-on-empty-changeset \
	 	--parameter-overrides \
	 		SlackSigningSecret=${SLACK_SIGNING_SECTET} \
	 		EnvironmentEncryptionKeyArn=${ENVIRONMENT_ENCRYPTION_KEY_ARN} \
	 		BotUserAccessToken=${BOT_USER_ACCESS_TOKEN} \
	 		DefaultChannelID=${DEFAULT_CHANNEL_ID}

show-api-url:
	 aws cloudformation describe-stacks \
	 	--stack-name slack-command-example2-application \
		--query "Stacks[0].Outputs[?OutputKey=='ProdDataEndpoint'].OutputValue | [0]" \
		--output text

sam-local-start: build
	[ -f build/sam.pid ] || (sam local start-api -t cfn/application.yaml & echo $$! > build/sam.pid)
	[ -f build/realize.pid ] || (GOARCH=amd64 GOOS=linux realize start & echo $$! > build/realize.pid)

sam-local-stop:
	[ -f build/sam.pid ] && kill `cat build/sam.pid` && rm build/sam.pid
	[ -f build/realize.pid ] && kill -9 `cat build/realize.pid` && rm build/realize.pid
