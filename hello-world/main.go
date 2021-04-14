package main

import (
	"encoding/json"
	"fmt"
	"hello-world/data"
	"hello-world/data/models"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-openapi/strfmt"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch request.HTTPMethod {
	case http.MethodGet:
		return get(request)
	case http.MethodPost:
		return post(request)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func notImplementedError() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotImplemented,
		Body:       http.StatusText(http.StatusNotImplemented),
	}
}

func get(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return notImplementedError(), nil
}

func post(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	input := new(models.OBReadConsent1)
	err := json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}

	err = input.Validate(strfmt.Default)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}

	response, err := json.Marshal(data.NameResponse{
		Message: fmt.Sprintf("Hello, %s %s", input.Data.ExpirationDateTime, input.Data.TransactionFromDateTime),
	})
	if err != nil {
		return serverError(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(response),
	}, nil
}

func main() {
	lambda.Start(handler)
}
