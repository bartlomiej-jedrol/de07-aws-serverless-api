package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/user"
	"github.com/stretchr/testify/assert"
)

var (
	validEmail   string = "bartlomiej.jedrol@gmail.com"
	invalidEmail string = "test.test@gmail.com"
	validBody    string = `{"email":"bartlomiej.jedrol@gmail.com","firstName":"Bartlomiej","lastName":"Jedrol","age":37}`
	invalidBody  string = `{"email":""bartlomiej.jedrol@gmail.com","firstName":"Bartlomiej","lastName":"Jedrol","age":37}`

	// Requests.
	requestValidUser = events.APIGatewayProxyRequest{
		HTTPMethod:            "GET",
		QueryStringParameters: map[string]string{"email": validEmail},
	}
	requestInvalidUser = events.APIGatewayProxyRequest{
		HTTPMethod:            "GET",
		QueryStringParameters: map[string]string{"email": invalidEmail},
	}
	requestEmptyUser = events.APIGatewayProxyRequest{
		HTTPMethod:            "GET",
		QueryStringParameters: map[string]string{"email": ""},
	}
	// requestNoQueryParameter = events.APIGatewayProxyRequest{
	// 	HTTPMethod:            "GET",
	// 	QueryStringParameters: map[string]string{},
	// }

	// Responses.
	responseStatusOK = events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Body:              validBody,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string(nil),
		IsBase64Encoded:   false,
	}
	responseBadRequest = events.APIGatewayProxyResponse{
		StatusCode:        http.StatusBadRequest,
		Body:              fmt.Sprintf(`{"error":"%v"}`, ErrorBadRequest.Error()), // "{\"error\":\"bad request\"}"
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string(nil),
		IsBase64Encoded:   false,
	}
	responseNotFound = events.APIGatewayProxyResponse{
		StatusCode:        http.StatusNotFound,
		Body:              fmt.Sprintf(`{"error":"%v"}`, ErrorNotFound.Error()),
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string(nil),
		IsBase64Encoded:   false,
	}
)

func TestUnmarshalUser(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   string
		expected      *user.User
		expectedError bool
	}{
		{
			name:        "Valid JSON",
			requestBody: validBody,
			expected: &user.User{
				Email:     validEmail,
				FirstName: "Bartlomiej",
				LastName:  "Jedrol",
				Age:       37,
			},
			expectedError: false,
		},
		{
			name:          "Invalid JSON",
			requestBody:   invalidBody,
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "Empty JSON",
			requestBody:   `{}`,
			expected:      &user.User{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unmarshalUser(tt.requestBody)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name     string
		request  *events.APIGatewayProxyRequest
		expected *events.APIGatewayProxyResponse
	}{
		{
			name:     "Valid user",
			request:  &requestValidUser,
			expected: &responseStatusOK,
		},
		{
			name:     "Invalid user",
			request:  &requestInvalidUser,
			expected: &responseNotFound,
		},
		{
			name:     "Empty user",
			request:  &requestEmptyUser,
			expected: &responseBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUser(*tt.request)
			assert.Equal(t, tt.expected, result)
		})
	}
}
