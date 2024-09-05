package user

import (
	"testing"

	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/models"
	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFetchUser(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:  "Valid user",
			email: testutil.ValidUser1.Email,
			expectedUser: &models.User{
				Email:     testutil.ValidUser1.Email,
				FirstName: testutil.ValidUser1.FirstName,
				LastName:  testutil.ValidUser1.LastName,
				Age:       testutil.ValidUser1.Age,
			},
			expectedError: nil,
		},
		{
			name:          "Invalid user",
			email:         testutil.InvalidUser1.Email,
			expectedUser:  nil,
			expectedError: ErrorUserDoesNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := FetchUser(tt.email)

			if tt.expectedError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, *tt.expectedUser, *user)
				}
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name          string
		user          models.User
		expectedError error
	}{
		{
			name: "Successful user creation",
			user: models.User{
				Email:     testutil.ValidUser1.Email,
				FirstName: testutil.ValidUser1.FirstName,
				LastName:  testutil.ValidUser1.LastName,
				Age:       testutil.ValidUser1.Age,
			},
			expectedError: nil,
		},
		{
			name: "Unsuccessful user creation",
			user: models.User{
				Email:     "",
				FirstName: testutil.ValidUser1.FirstName,
				LastName:  testutil.ValidUser1.LastName,
				Age:       testutil.ValidUser1.Age,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateUser(tt.user)

			if tt.expectedError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.expectedError, err)
				}
			}
			// else {
			// 	if assert.NoError(t, err) {
			// 		assert.Equal(t, *tt.expectedUser, *user)
			// 	}
			// }
		})
	}
}
