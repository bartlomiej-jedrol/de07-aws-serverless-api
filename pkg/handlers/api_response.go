// Handlers implements handlers for all possible HTTP methods and the API response.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func buildAPIResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v", err)
		return nil, err
	}

	responseBody := &events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: status,
		Body:       string(bodyJSON),
	}
	log.Printf("========== responseBody ==========: %v", responseBody)
	return responseBody, nil
}
