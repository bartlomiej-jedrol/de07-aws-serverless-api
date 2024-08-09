// Handlers implements handlers for all possible HTTP methods and the API response.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"

	"github.com/aws/aws-lambda-go/events"
)

// buildResponseBody returns body of the API response based on the status code.
func buildResponseBody(status int, body interface{}) string {
	successfulStatuses := []int{200, 201}
	// If the status is not successful then add error message as body.
	if !slices.Contains(successfulStatuses, status) {
		// Convert error to string if exists.
		var errorMessage string
		if err, ok := body.(error); ok {
			errorMessage = err.Error()
		} else {
			errorMessage = fmt.Sprintf("%v", body)
		}

		body := map[string]interface{}{"error": errorMessage}
		responseBody, err := json.Marshal(body)
		if err != nil {
			fmt.Printf("Failed to marshal JSON: %v", err)
		}
		return string(responseBody)
		// Else keep body as is.
	} else {
		responseBody, err := json.Marshal(body)
		if err != nil {
			fmt.Printf("Failed to marshal JSON: %v", err)
		}
		return string(responseBody)
	}
}

// buildAPIResponse builds API response.
func buildAPIResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	// Build response body.
	responseBody := &events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: status,
	}

	responseBody.Body = buildResponseBody(status, body)
	log.Printf("========== responseBody ==========: %v", responseBody)
	return responseBody, nil
}
