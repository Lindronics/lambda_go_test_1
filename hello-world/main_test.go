package main

import (
	"encoding/json"
	"hello-world/data"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"gopkg.in/go-playground/validator.v9"
)

func TestHandler(t *testing.T) {

	t.Run("Invalid Request JSON", func(t *testing.T) {
		response, err := handler(events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
		})
		if err != nil {
			t.Fatal("Error")
		}
		if response.StatusCode != http.StatusBadRequest {
			t.Fatal("Response code is not invalid request")
		}
	})

	t.Run("Missing fields", func(t *testing.T) {
		response, err := handler(events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Body: `
				{
					"first_name": "asdf"
				}
			`,
		})
		if err != nil {
			t.Fatal("Error")
		}
		if response.StatusCode != http.StatusBadRequest {
			t.Fatal("Response code is not invalid request")
		}
	})

	t.Run("Successful request", func(t *testing.T) {
		response, err := handler(events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Body: `
				{
					"first_name": "asdf",
					"last_name": "asdf"
				}
			`,
		})
		if err != nil {
			t.Fatal("Error")
		}
		if response.StatusCode != http.StatusCreated {
			t.Fatal("Response code is not success")
		}
		output := new(data.NameResponse)
		err = json.Unmarshal([]byte(response.Body), &output)
		if err != nil {
			t.Fatal("Response not valid JSON")
		}

		validate := validator.New()
		err = validate.Struct(output)
		if err != nil {
			t.Fatal("Invalid JSON schema")
		}
	})
}
