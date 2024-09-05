package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/models"
	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/testutil"
	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/user"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalUser(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   string
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:        "Valid JSON",
			requestBody: testutil.ValidUser,
			expectedUser: &models.User{
				Email:     testutil.ValidUser1.Email,
				FirstName: testutil.ValidUser1.FirstName,
				LastName:  testutil.ValidUser2.LastName,
				Age:       testutil.ValidUser1.Age,
			},
			expectedError: nil,
		},
		{
			name:          "Invalid User",
			requestBody:   testutil.InvalidJSON,
			expectedUser:  nil,
			expectedError: ErrorInvalidJSON,
		},
		{
			name:          "Empty User",
			requestBody:   testutil.EmptyUser,
			expectedUser:  &models.User{},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualUser, err := unmarshalUser(tt.requestBody)

			if tt.expectedError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, *tt.expectedUser, *actualUser)
				}
			}
		})
	}
}

func TestMapErrorToResponse(t *testing.T) {
	tests := []struct {
		name               string
		inputError         error
		expectedStatusCode int
		expectedError      error
	}{
		{
			name:               "ErrorFailedToGetItems",
			inputError:         user.ErrorFailedToGetItems,
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      ErrorInternalServerError,
		},
		{
			name:               "ErrorFailedToPutItem",
			inputError:         user.ErrorFailedToPutItem,
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      ErrorInternalServerError,
		},
		{
			name:               "ErrorFailedToDeleteItem",
			inputError:         user.ErrorFailedToDeleteItem,
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      ErrorInternalServerError,
		},
		{
			name:               "ErrorFailedToUnmarshalMap",
			inputError:         user.ErrorFailedToUnmarshalMap,
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      ErrorInternalServerError,
		},
		{
			name:               "ErrorUserDoesNotExist",
			inputError:         user.ErrorUserDoesNotExist,
			expectedStatusCode: http.StatusNotFound,
			expectedError:      ErrorNotFound,
		},
		{
			name:               "ErrorFailedToValidateUser",
			inputError:         user.ErrorFailedToValidateUser,
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      ErrorBadRequest,
		},
		{
			name:               "ErrorInvalidJSON",
			inputError:         ErrorInvalidJSON,
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      ErrorBadRequest,
		},
		{
			name:               "UnknownError",
			inputError:         errors.New("unknown error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      ErrorInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, err := mapErrorToResponse(tt.inputError)
			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Equal(t, tt.expectedError, err)
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
				QueryStringParameters: testutil.ValidQueQueryStringParameters,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       testutil.ValidUser,
			},
		},
		{
			name: "Invalid user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            "GET",
				QueryStringParameters: testutil.InvalidQueQueryStringParameters,
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
			assert.Equal(t, tt.expected.StatusCode, actual.StatusCode)
			assert.JSONEq(t, tt.expected.Body, actual.Body)
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
				Body:       testutil.ValidUsers,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := GetUsers()
			assert.Equal(t, tt.expected.StatusCode, actual.StatusCode)
			assert.JSONEq(t, tt.expected.Body, actual.Body)
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
				Body:       testutil.ValidUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusCreated,
				Body:       testutil.ValidUser,
			},
		},
		{
			name: "Empty user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       testutil.UserEmptyEmail,
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
				Body:       testutil.InvalidJSON,
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
			assert.Equal(t, tt.expected.StatusCode, actual.StatusCode)
			assert.JSONEq(t, tt.expected.Body, actual.Body)
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
				Body:       testutil.ValidUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       testutil.ValidUser,
			},
		},
		{
			name: "Invalid user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "PUT",
				Body:       testutil.InvalidUser,
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
				Body:       testutil.UserEmptyEmail,
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
				Body:       testutil.InvalidJSON,
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
			assert.Equal(t, tt.expected.StatusCode, actual.StatusCode)
			assert.JSONEq(t, tt.expected.Body, actual.Body)
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
				Body:       testutil.ValidUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       testutil.ValidUser,
			},
		},
		{
			name: "Invalid user",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "DELETE",
				Body:       testutil.InvalidUser,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       fmt.Sprintf(`{"error":"%v"}`, ErrorNotFound.Error()),
			},
		},
		{
			name: "User empty email",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "DELETE",
				Body:       testutil.EmptyUser,
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
				Body:       testutil.InvalidJSON,
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
			assert.Equal(t, tt.expected.StatusCode, actual.StatusCode)
			assert.JSONEq(t, tt.expected.Body, actual.Body)
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
				Body:       testutil.ValidUser,
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
			assert.Equal(t, tt.expected.StatusCode, actual.StatusCode)
			assert.JSONEq(t, tt.expected.Body, actual.Body)
		})
	}
}
