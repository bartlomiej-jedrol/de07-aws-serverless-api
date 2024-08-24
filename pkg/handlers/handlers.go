// Handlers implements handlers for all possible HTTP methods and the API response.
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/user"
)

var (
	// ErrorMethodNotSupported Not supported method
	ErrorMethodNotAllowed = errors.New("method not supported")

	// ErrorBadRequest Bad request
	ErrorBadRequest = errors.New("bad request")

	// ErrorNotFound Not found
	ErrorNotFound = errors.New("not found")
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

// GetUser gets the user data from the DynamoDB table.
// It also returns the response to the caller.
func GetUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Extract users's email from the request.
	email := request.QueryStringParameters["email"]

	// Fetch user data from DynamoDB. If err return not found.
	userData, err := user.FetchUser(email)
	if err != nil {
		return buildAPIResponse(http.StatusNotFound, ErrorNotFound)
	}
	return buildAPIResponse(http.StatusOK, userData)
}

// CreateUser extracts user data from the request and creates the user in the DynamoDB table.
// It also returns the response to the caller.
func CreateUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Unmarshal received user JSON data.
	userJSON := request.Body
	var userData user.User
	err := json.Unmarshal([]byte(userJSON), &userData)
	if err != nil {
		log.Printf("Failed to unmarshal userJSON: %v", err)
	}
	log.Printf("========== User ==========: %v", userData)

	// Create user.
	err = user.CreateUser(userData)
	if err != nil {
		return buildAPIResponse(http.StatusBadRequest, ErrorBadRequest)
	}

	// Send response for successful user creation.
	return buildAPIResponse(http.StatusCreated, userData)
}

func UpdateUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var userData user.User
	userJSON := request.Body
	err := json.Unmarshal([]byte(userJSON), &userData)
	if err != nil {
		log.Printf("Failed to unmarshal userJSON: %v", err)
	}
	user.UpdateUser(userData)
	return buildAPIResponse(http.StatusOK, userData)
}

func DeleteUser() {

}

func UnhandledHTTPMethod(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Println("UnhandledHTTPMethod called")
	return buildAPIResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
