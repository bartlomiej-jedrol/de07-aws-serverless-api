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

	// ErrorInvalidJSON Invalid JSON
	ErrorInvalidJSON = errors.New("invalid JSON")

	// ErrorNoEmailQueryParameter No email query parameter
	ErrorNoEmailQueryParameter = errors.New("no email query parameter")

	// ErrorBadRequest Bad request
	ErrorBadRequest = errors.New("bad request")

	// ErrorNotFound Not found
	ErrorNotFound = errors.New("not found")

	// ErrorInternalServerError Internal server error
	ErrorInternalServerError = errors.New("internal server error")
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func unmarshalUser(request events.APIGatewayProxyRequest) (*user.User, error) {
	var u user.User
	err := json.Unmarshal([]byte(request.Body), &u)
	if err != nil {
		log.Printf("%v: %v", ErrorInvalidJSON, err)
		return nil, ErrorInvalidJSON
	}
	log.Printf("User: %v", u)

	return &u, nil
}

func mapErrorToResponse(err error) (int, error) {
	switch err {
	case user.ErrorFailedToGetItem, user.ErrorFailedToPutItem, user.ErrorFailedToDeleteItem, user.ErrorFailedToUnmarshalMap:
		return http.StatusInternalServerError, ErrorInternalServerError
	case user.ErrorUserDoesNotExist:
		return http.StatusNotFound, ErrorNotFound
	case user.ErrorFailedToValidateUser, ErrorInvalidJSON:
		return http.StatusBadRequest, ErrorBadRequest
	default:
		return http.StatusInternalServerError, ErrorInternalServerError
	}
}

// GetUser gets the user data from the DynamoDB table.
// It also returns the response to the caller.
func GetUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Extract users's email from the request.
	email := request.QueryStringParameters["email"]
	if email == "" {
		log.Printf("%v", ErrorNoEmailQueryParameter)
		return buildAPIResponse(http.StatusBadRequest, ErrorNoEmailQueryParameter)
	}
	log.Printf("query parameter email: %v", email)

	// Fetch user.
	u, err := user.FetchUser(email)
	if err != nil {
		statusCode, errorMessage := mapErrorToResponse(err)
		return buildAPIResponse(statusCode, errorMessage)
	}

	// Send successful response.
	return buildAPIResponse(http.StatusOK, u)
}

// CreateUser extracts user data from the request and creates the user in the DynamoDB table.
// It also returns the response to the caller.
func CreateUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Unmarshal received user JSON data.
	u, err := unmarshalUser(request)
	if err != nil {
		statusCode, errorMessage := mapErrorToResponse(err)
		return buildAPIResponse(statusCode, errorMessage)
	}

	// Create user.
	err = user.CreateUser(*u)
	if err != nil {
		statusCode, errorMessage := mapErrorToResponse(err)
		return buildAPIResponse(statusCode, errorMessage)
	}

	// Send successful response.
	return buildAPIResponse(http.StatusCreated, u)
}

func UpdateUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Unmarshal received user JSON data.
	u, err := unmarshalUser(request)
	if err != nil {
		statusCode, errorMessage := mapErrorToResponse(err)
		return buildAPIResponse(statusCode, errorMessage)
	}

	// Update user.
	err = user.UpdateUser(*u)
	if err != nil {
		statusCode, errorMessage := mapErrorToResponse(err)
		return buildAPIResponse(statusCode, errorMessage)
	}

	// Send successful response.
	return buildAPIResponse(http.StatusOK, u)
}

func DeleteUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Extract users's email from the request.
	email := request.QueryStringParameters["email"]
	if email == "" {
		return buildAPIResponse(http.StatusBadRequest, ErrorNoEmailQueryParameter)
	}
	log.Printf("query parameter email: %v", email)

	// Delete item from DynamoDB table.
	err := user.DeleteUser(email)
	if err != nil {
		statusCode, errorMessage := mapErrorToResponse(err)
		return buildAPIResponse(statusCode, errorMessage)
	}

	// Send successful response.
	return buildAPIResponse(http.StatusOK, request.Body)
}

func UnhandledHTTPMethod(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Println("UnhandledHTTPMethod called")
	return buildAPIResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
