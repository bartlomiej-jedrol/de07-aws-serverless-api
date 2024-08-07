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

	// ErrorMethodNotSupported Bad request
	ErrorBadRequest = errors.New("bad request")
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

// func GetUser(request events.APIGatewayProxyRequest, dynamodbClient dynamodb.Client, tableName string) *events.APIGatewayProxyResponse {
// 	fetchedUser := user.FetchUser()
// }

// CreateUser extracts user data from the request and passes to the user.CreateUser function that creates the user in the DynamoDB table.
// It also returns the response to the caller.
func CreateUser(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Unmarshal received user JSON data.
	userJSON := request.Body
	var userData user.User
	err := json.Unmarshal([]byte(userJSON), &userData)
	if err != nil {
		log.Fatalf("Failed to unmarshal userJSON: %v", err)
	}
	log.Printf("========== User ==========: %v", userData)

	// Create user.
	if err := user.CreateUser(userData); err != nil {
		return buildAPIResponse(http.StatusBadRequest, ErrorBadRequest)
	}

	// Send response for successful user creation.
	return buildAPIResponse(http.StatusCreated, userData)
}

func UpdateUser() {

}

func DeleteUser() {

}

func UnhandledHTTPMethod(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Println("UnhandledHTTPMethod called")
	return buildAPIResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
