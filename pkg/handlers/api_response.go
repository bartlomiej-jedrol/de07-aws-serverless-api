// Handlers implements handlers for all possible HTTP methods and the API response.
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/aws/aws-lambda-go/events"
)

var (
	ErrorFailedToMarshalJSON = errors.New("failed to marshal JSON")
)

// buildResponseBody returns body of the API response based on the status code.
func buildResponseBody(status int, body interface{}) (int, string) {
	successfulStatuses := []int{200, 201}
	// If the status is not successful then add error message to response body.
	if !slices.Contains(successfulStatuses, status) {
		// If exists, convert error to string.
		var errorMessage string
		if err, ok := body.(error); ok {
			errorMessage = err.Error()
		} else {
			errorMessage = fmt.Sprintf("%v", body)
		}

		// Build body with error message.
		body := map[string]interface{}{"error": errorMessage}
		responseBody, err := json.Marshal(body)
		if err != nil {
			log.Printf("%v: %v", ErrorFailedToMarshalJSON, err)
			return 500, ""
		}
		return status, string(responseBody)
		// Else keep body as is.
	} else {
		responseBody, err := json.Marshal(body)
		if err != nil {
			log.Printf("%v: %v", ErrorFailedToMarshalJSON, err)
			return 500, ""
		}
		return status, string(responseBody)
	}
}

// buildAPIResponse builds API response.
func buildAPIResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	// Build response body.
	responseBody := &events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: status,
	}

	responseBody.StatusCode, responseBody.Body = buildResponseBody(status, body)
	log.Printf("responseBody: %v", responseBody)
	return responseBody, nil
}
