package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"hello-world/data"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/go-playground/validator.v9"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
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
	// errorLogger.Println(err.Error())

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

	input := new(data.Name)
	err := json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       fmt.Sprintf("{\"message\": \"Hello, %s %s\"}", input.FirstName, input.LastName),
	}, nil
}

func main() {
	lambda.Start(handler)
}
