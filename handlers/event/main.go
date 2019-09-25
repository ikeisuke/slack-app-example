package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ikeisuke/slack-app-example/internal/application"
	"github.com/ikeisuke/slack-app-example/internal/infrastructure"
	"github.com/ikeisuke/slack-app-example/internal/presenter"
	"github.com/ikeisuke/slack-app-example/internal/repository"
	"os"
	"strconv"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error
	timestamp, err := strconv.Atoi(request.Headers["X-Slack-Request-Timestamp"])
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	repo := repository.NewSignatureRepository()
	present := presenter.NewSimplePresenter()
	infra := infrastructure.NewSlack(os.Getenv("BOT_USER_ACCESS_TOKEN"), os.Getenv("DEFAULT_CHANNEL_ID義tq"))
	event := repository.NewEventRepository()
	app := application.NewEventReceiverInteraction(repo, event, present, infra)
	res := app.Run(&application.EventReceiverInput{
		Timestamp:        timestamp,
		Signature:        request.Headers["X-Slack-Signature"],
		SigningSecret:    os.Getenv("SLACK_SIGNING_SECRET"),
		Body:             request.Body,
		SignatureVersion: "v0",
	})
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            res,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
