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
	fmt.Printf("%+v", request)
	repo := repository.NewSignatureRepository()
	present := presenter.NewResponsePresenter()
	event := repository.NewEventRepository()
	app := application.NewEventReceiverInteraction(repo, event, present)
	res := app.Run(&application.EventReceiverInput{
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
	//type urlVerificationEvent struct {
	//	Token     string `json:"token"`
	//	Challenge string `json:"challenge"`
	//	Type      string `json:"type"`
	//}
	//e := &urlVerificationEvent{}
	//if err := json.Unmarshal([]byte(request.Body), e); err != nil {
	//	fmt.Errorf("%+v", err)
	//	os.Exit(1)
	//}
	//body, _ := json.Marshal(map[string]string{
	//	"challenge": e.Challenge,
	//})
	//fmt.Printf("%+v", request)
	//return events.APIGatewayProxyResponse{
	//	StatusCode:      200,
	//	IsBase64Encoded: false,
	//	Body:            string(body),
	//	Headers: map[string]string{
	//		"Content-Type": "application/json",
	//	},
	//}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
