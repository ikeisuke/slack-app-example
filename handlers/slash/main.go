package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ikeisuke/slack-app-example/internal/application"
	"github.com/ikeisuke/slack-app-example/internal/presenter"
	"github.com/ikeisuke/slack-app-example/internal/repository"
	"os"
	"strconv"
)

func HandleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error
	timestamp, err := strconv.Atoi(request.Headers["X-Slack-Request-Timestamp"])
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	repo := repository.NewSignatureRepository()
	present := presenter.NewResponsePresenter()
	sub := repository.NewSubCommandRepository()
	app := application.NewSlashCommandInteraction(repo, sub, present)
	res := app.Run(&application.SlashCommandInput{
		Timestamp:        timestamp,
		Signature:        request.Headers["X-Slack-Signature"],
		SigningSecret:    os.Getenv("SLACK_SIGNING_SECRET"),
		Body:             request.Body,
		SignatureVersion: "v0",
	})
	fmt.Printf("%+v", res)
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
