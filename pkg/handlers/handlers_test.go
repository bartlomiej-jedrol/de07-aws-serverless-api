package handlers

import (
	"errors"
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
	validUser    string = `{"email":"bartlomiej.jedrol@gmail.com","firstName":"Bartlomiej","lastName":"Jedrol","age":37}`
	invalidUser  string = `{"email":"test.test@gmail.com","firstName":"Bartlomiej","lastName":"Jedrol","age":37}`
	emptyUser    string = `{"email":"","firstName":"Bartlomiej","lastName":"Jedrol","age":37}`
	invalidJSON  string = `{"email":""bartlomiej.jedrol@gmail.com","firstName":"Bartlomiej","lastName":"Jedrol","age":37}`
	validUsers   string = `[{"email":"jedrol.natalia@gmail.com","firstName":"Natalia","lastName":"Jedrol","age":33},{"email":"bartlomiej.jedrol@gmail.com","firstName":"Bartlomiej","lastName":"Jedrol","age":37}]`
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
			requestBody: validUser,
			expected: &user.User{
				Email:     validEmail,
				FirstName: "Bartlomiej",
				LastName:  "Jedrol",
				Age:       37,
			},
			expectedError: false,
		},
		{
			name:          "Invalid User",
			requestBody:   invalidJSON,
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "Empty User",
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

func TestMapErrorToResponse(t *testing.T) {
	tests := []struct {
		name               string
		inputErr           error
		expectedStatusCode int
		expectedErr        error
	}{
		{
			name:               "ErrorFailedToGetItems",
			inputErr:           user.ErrorFailedToGetItems,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        ErrorInternalServerError,
		},
		{
			name:               "ErrorFailedToPutItem",
			inputErr:           user.ErrorFailedToPutItem,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        ErrorInternalServerError,
		},
		{
			name:               "ErrorFailedToDeleteItem",
			inputErr:           user.ErrorFailedToDeleteItem,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        ErrorInternalServerError,
		},
		{
			name:               "ErrorFailedToUnmarshalMap",
			inputErr:           user.ErrorFailedToUnmarshalMap,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        ErrorInternalServerError,
		},
		{
			name:               "ErrorUserDoesNotExist",
			inputErr:           user.ErrorUserDoesNotExist,
			expectedStatusCode: http.StatusNotFound,
			expectedErr:        ErrorNotFound,
		},
		{
			name:               "ErrorFailedToValidateUser",
			inputErr:           user.ErrorFailedToValidateUser,
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        ErrorBadRequest,
		},
		{
			name:               "ErrorInvalidJSON",
			inputErr:           ErrorInvalidJSON,
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        ErrorBadRequest,
		},
		{
			name:               "UnknownError",
			inputErr:           errors.New("unknown error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        ErrorInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, err := mapErrorToResponse(tt.inputErr)
			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name     string
		request  events.APIGatewayProxyRequest
		expected events.APIGatewayProxyResponse
	}{
		{
			name: "Valid user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            "GET",
				QueryStringParameters: map[string]string{"email": validEmail},
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       validUser,
			},
		},
		{
			name: "Invalid user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            "GET",
				QueryStringParameters: map[string]string{"email": invalidEmail},
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorNotFound.Error()),
			},
		},
		{
			name: "Empty user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            "GET",
				QueryStringParameters: map[string]string{"email": ""},
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorBadRequest.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := GetUser(tt.request)
			assert.Equal(t, actual.StatusCode, tt.expected.StatusCode)
			assert.JSONEq(t, actual.Body, tt.expected.Body)
		})
	}
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name     string
		expected events.APIGatewayProxyResponse
	}{
		{
			name: "Successful retrieval of users",
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       validUsers,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := GetUsers()
			assert.Equal(t, actual.StatusCode, tt.expected.StatusCode)
			assert.JSONEq(t, actual.Body, tt.expected.Body)
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name     string
		request  events.APIGatewayProxyRequest
		expected events.APIGatewayProxyResponse
	}{
		{
			name: "Valid user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       validUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusCreated,
				Body:       validUser,
			},
		},
		{
			name: "Empty user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       emptyUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorBadRequest.Error()),
			},
		},
		{
			name: "Invalid JSON",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       invalidJSON,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorBadRequest.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := CreateUser(tt.request)
			t.Log(actual)
			assert.Equal(t, actual.StatusCode, tt.expected.StatusCode)
			assert.JSONEq(t, actual.Body, tt.expected.Body)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name     string
		request  events.APIGatewayProxyRequest
		expected events.APIGatewayProxyResponse
	}{
		{
			name: "Valid user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "PUT",
				Body:       validUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       validUser,
			},
		},
		{
			name: "Invalid user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "PUT",
				Body:       invalidUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorNotFound.Error()),
			},
		},
		{
			name: "Empty user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "PUT",
				Body:       emptyUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorBadRequest.Error()),
			},
		},
		{
			name: "Invalid JSON",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "PUT",
				Body:       invalidJSON,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorBadRequest.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("tt: %v", tt)
			actual, _ := UpdateUser(tt.request)
			t.Logf("actual: %v", actual)
			assert.Equal(t, actual.StatusCode, tt.expected.StatusCode)
			assert.JSONEq(t, actual.Body, tt.expected.Body)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name     string
		request  events.APIGatewayProxyRequest
		expected events.APIGatewayProxyResponse
	}{
		{
			name: "Valid user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "DELETE",
				Body:       validUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       validUser,
			},
		},
		{
			name: "Invalid user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "DELETE",
				Body:       invalidUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorNotFound.Error()),
			},
		},
		{
			name: "Empty user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "DELETE",
				Body:       emptyUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorBadRequest.Error()),
			},
		},
		{
			name: "Invalid JSON",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "DELETE",
				Body:       invalidJSON,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorBadRequest.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("tt: %v", tt)
			actual, _ := UpdateUser(tt.request)
			t.Logf("actual: %v", actual)
			assert.Equal(t, actual.StatusCode, tt.expected.StatusCode)
			assert.JSONEq(t, actual.Body, tt.expected.Body)
		})
	}
}

func TestUnhandledHTTPMethod(t *testing.T) {
	tests := []struct {
		name     string
		request  events.APIGatewayProxyRequest
		expected events.APIGatewayProxyResponse
	}{
		{
			name: "Unhandled HTTP method",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "PATCH",
				Body:       validUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusMethodNotAllowed,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorMethodNotAllowed.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("tt: %v", tt)
			actual, _ := UnhandledHTTPMethod(tt.request)
			t.Logf("actual: %v", actual)
			assert.Equal(t, actual.StatusCode, tt.expected.StatusCode)
			assert.JSONEq(t, actual.Body, tt.expected.Body)
		})
	}
}
