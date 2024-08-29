package handlers

import (
	"testing"

	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/user"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalUser(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   string
		expected      *user.User
		expectedError bool
	}{
		{
			name: "Valid JSON",
			requestBody: `{"email": "bartlomiej.jedrol@gmail.com", "firstName": "Bartlomiej", 
			"lastName": "Jedrol", "age": 37}`,
			expected: &user.User{
				Email:     "bartlomiej.jedrol@gmail.com",
				FirstName: "Bartlomiej",
				LastName:  "Jedrol",
				Age:       37,
			},
			expectedError: false,
		},
		{
			name: "Invalid JSON",
			requestBody: `{""email": "bartlomiej.jedrol@gmail.com", "firstName": "Bartlomiej", 
			"lastName": "Jedrol", "age": 37}`,
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
