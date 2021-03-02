package main

import (
	"encoding/json"
	"fmt"
	"hello-world/data"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"gopkg.in/go-playground/validator.v9"
)

func TestMain(m *testing.M) {
	rc := m.Run()

	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.8 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}

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
