package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/ezzarghili/recaptcha-go.v4"
)

func main() {
	runtime.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{StatusCode: 405}, nil

	}

	captcha, _ := recaptcha.NewReCAPTCHA(os.Getenv("RECAPTCHA_SECRET"), recaptcha.V2, 10*time.Second) // for v2 API get your secret from https://www.google.com/recaptcha/admin

	var body map[string]interface{}

	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		// only allow post requests
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	if token, ok := body["token"]; ok {
		if err := captcha.Verify(token.(string)); err != nil {
			// recaptcha validation failed
			return events.APIGatewayProxyResponse{StatusCode: 403}, nil
		}
	} else {
		// recaptcha token not provided
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	if site, ok := body["site"]; ok {
		email := os.Getenv(site.(string))
		if email == "" {
			// site provided was not found
			return events.APIGatewayProxyResponse{StatusCode: 403}, nil
		}

		// all is good, return the email address
		return events.APIGatewayProxyResponse{Body: "{\"email\":\"" + email + "\"}", StatusCode: 200}, nil
	} else {
		// site not provided
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}
}
