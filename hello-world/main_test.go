package main

import (
	"encoding/json"
	"errors"
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
		if c < 0.7 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}

func TestHandler(t *testing.T) {
	t.Run("HTTP Request methods check", func(t *testing.T) {
		response, err := handler(events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodDelete,
		})
		if err != nil {
			t.Fatal("Error")
		}
		if response.StatusCode != http.StatusMethodNotAllowed {
			t.Fatal("Should return method not allowed")
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("Status check", func(t *testing.T) {
		response, err := handler(events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodGet,
		})
		if err != nil {
			t.Fatal("Error")
		}
		if response.StatusCode != http.StatusNotImplemented {
			t.Fatal("Should return not implemented")
		}
	})
}

func TestPost(t *testing.T) {

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
				"Data": {
					"ExpirationDateTime": "2021-04-13T21:36:16.318Z",
					"Permissions": [
					"ReadAccountsBasic"
					],
					"TransactionFromDateTime": "2021-04-13T21:36:16.318Z",
					"TransactionToDateTime": "2021-04-13T21:36:16.318Z"
				},
				"Risk": {}
			}`,
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

func TestServerError(t *testing.T) {
	t.Run("Server error check", func(t *testing.T) {
		response, err := serverError(errors.New("test error"))
		if err != nil {
			t.Fatal("Error")
		}
		if response.StatusCode != http.StatusInternalServerError {
			t.Fatal("HTTP status should be 500")
		}
	})
}
