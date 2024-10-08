// Handlers implements handlers for all possible HTTP methods and an API response.
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/models"
	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/user"
)

var (
	ErrorMethodNotAllowed      = errors.New("method not supported")
	ErrorInvalidJSON           = errors.New("invalid JSON")
	ErrorNoEmailQueryParameter = errors.New("no email query parameter")
	ErrorBadRequest            = errors.New("bad request")
	ErrorNotFound              = errors.New("not found")
	ErrorInternalServerError   = errors.New("internal server error")
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

// unmarshalUser unmarshals user from body.
func unmarshalUser(body string) (*models.User, error) {
	var u models.User
	err := json.Unmarshal([]byte(body), &u)
	if err != nil {
		log.Printf("%v: %v", ErrorInvalidJSON, err)
		return nil, ErrorInvalidJSON
	}
	log.Printf("User: %v", u)

	return &u, nil
}

// mapErrorToResponse maps business logic errors to the HTTP response errors and status codes.
func mapErrorToResponse(err error) (int, error) {
	switch err {
	case user.ErrorFailedToGetItem, user.ErrorFailedToGetItems, user.ErrorFailedToPutItem,
		user.ErrorFailedToDeleteItem, user.ErrorFailedToUnmarshalMap:
		return http.StatusInternalServerError, ErrorInternalServerError
	case user.ErrorUserDoesNotExist:
		return http.StatusNotFound, ErrorNotFound
	case user.ErrorFailedToValidateUser, ErrorInvalidJSON:
		return http.StatusBadRequest, ErrorBadRequest
	default:
		return http.StatusInternalServerError, ErrorInternalServerError
	}
}

// GetUser gets user data from DynamoDB and responds.
func GetUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Extract users's email from request.
	email := request.QueryStringParameters["email"]
	if email == "" {
		log.Printf("%v", ErrorNoEmailQueryParameter)
		return buildAPIResponse(http.StatusBadRequest, ErrorBadRequest)
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

// GetUsers gets users' data from DynamoDB table and responds.
func GetUsers() (*events.APIGatewayProxyResponse, error) {
	// Fetch users.
	users, err := user.FetchUsers()
	if err != nil {
		statusCode, errorMessage := mapErrorToResponse(err)
		return buildAPIResponse(statusCode, errorMessage)
	}

	// Send successful response.
	return buildAPIResponse(http.StatusOK, users)
}

// CreateUser creates user in DynamoDB table and responds.
func CreateUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Unmarshal received user JSON data.
	u, err := unmarshalUser(request.Body)
	if err != nil {
		statusCode, errorMessage := mapErrorToResponse(err)
		return buildAPIResponse(statusCode, errorMessage)
	}

	// Check email existence.
	if u.Email == "" {
		return buildAPIResponse(http.StatusBadRequest, ErrorBadRequest)
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

// UpdateUser updates user data in DynamoDB table and responds.
func UpdateUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Unmarshal received user JSON data.
	u, err := unmarshalUser(request.Body)
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

// DeleteUser deletes user data from DynamoDB table and responds.
func DeleteUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Extract users's email from request.
	email := request.QueryStringParameters["email"]
	if email == "" {
		log.Printf("%v", ErrorNoEmailQueryParameter)
		return buildAPIResponse(http.StatusBadRequest, ErrorBadRequest)
	}
	log.Printf("query parameter email: %v", email)

	// Delete item from DynamoDB table.
	u, err := user.DeleteUser(email)
	if err != nil {
		statusCode, errorMessage := mapErrorToResponse(err)
		return buildAPIResponse(statusCode, errorMessage)
	}

	// Send successful response.
	return buildAPIResponse(http.StatusOK, *u)
}

// UnhandledHTTPMethod responds for unsupported HTTP methods.
func UnhandledHTTPMethod(
	request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("unsupported HTTP method")
	return buildAPIResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
