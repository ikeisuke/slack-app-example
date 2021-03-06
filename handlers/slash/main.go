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
	timestamp, err := strconv.Atoi(request.Headers["X-Slack-Request-Timestamp"])
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	repo := repository.NewSignatureRepository()
	present := presenter.NewJSONPresenter()
	infra := infrastructure.NewSlack(os.Getenv("BOT_USER_ACCESS_TOKEN"), os.Getenv("DEFAULT_CHANNEL_ID"))
	sub := repository.NewCommandRepository(infra)
	app := application.NewSlashCommandInteraction(repo, sub, present)
	res, err := app.Run(&application.SlashCommandInput{
		Timestamp:        timestamp,
		Signature:        request.Headers["X-Slack-Signature"],
		SigningSecret:    os.Getenv("SLACK_SIGNING_SECRET"),
		Body:             request.Body,
		SignatureVersion: "v0",
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       res,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
