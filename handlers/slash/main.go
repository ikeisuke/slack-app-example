package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"strings"

	//"github.com/nlopes/slack"
)

func MakeHMAC(msg, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	timestamp := request.Headers["X-Slack-Request-Timestamp"]
	signature := request.Headers["X-Slack-Signature"]
	secret := os.Getenv("SLACK_SIGNING_SECRET")
	version := "v0"
	body := request.Body
	base := strings.Join([]string{
		version,
		timestamp,
		body,
	}, ":")
	sign := strings.Join([]string{
		version,
		MakeHMAC(base, secret),
	},"=")
	if signature != sign {
		fmt.Errorf("%s", "Invalid signature detected")
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			IsBase64Encoded: false,
			Body:            "",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	}
	tmp := make(map[string]string)
	for _, line := range strings.Split(body, "&") {
		kv := strings.Split(line, "=")
		if len(kv) != 2 {
			fmt.Errorf("%s", "Invalid request body detected")
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				IsBase64Encoded: false,
				Body:            body,
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}
		key := kv[0]
		value := kv[1]
		tmp[key] = value
	}
	buf, _ := json.Marshal(tmp)
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(buf),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
