package user

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func TestFetchUsers(t *testing.T) {
	tests := []struct {
		name          string
		expectedUsers []models.User
		expectedError error
	}{
		{
			name:          "Successful users retrieval",
			expectedUsers: []models.User{testutil.ValidUser2, testutil.ValidUser1},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualUsers, err := FetchUsers()

			if tt.expectedError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.expectedUsers, actualUsers)
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
			expectedError: ErrorFailedToPutItem,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateUser(tt.user)
			t.Logf("err:%v", err)

			if tt.expectedError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		user          models.User
		expectedError error
	}{
		{
			name:          "Valid user",
			user:          testutil.ValidUser1,
			expectedError: nil,
		},
		{
			name:          "Invalid user",
			user:          testutil.InvalidUser1,
			expectedError: ErrorUserDoesNotExist,
		},
		{
			name:          "Empty user",
			user:          testutil.EmptyUser1,
			expectedError: ErrorFailedToValidateUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateUser(tt.user)
			t.Logf("err:%v", err)

			if tt.expectedError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteUer(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		expectedUser  models.User
		expectedError error
	}{
		{
			name:          "Valid user",
			email:         testutil.ValidUser1.Email,
			expectedUser:  testutil.ValidUser1,
			expectedError: nil,
		},
		{
			name:          "Invalid user",
			email:         testutil.InvalidUser1.Email,
			expectedUser:  testutil.InvalidUser1,
			expectedError: ErrorUserDoesNotExist,
		},
		{
			name:          "Empty user",
			email:         testutil.EmptyUser1.Email,
			expectedUser:  testutil.InvalidUser1,
			expectedError: ErrorFailedToGetItem,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualUser, err := DeleteUser(tt.email)
			t.Logf("err:%v", err)

			if tt.expectedError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.expectedUser, *actualUser)
				}
			}
		})
	}
}

func TestGetKey(t *testing.T) {
	tests := []struct {
		name              string
		user              models.User
		expectedAttribute map[string]types.AttributeValue
	}{
		{
			name: "Valid user",
			user: testutil.ValidUser1,
			expectedAttribute: map[string]types.AttributeValue{
				"email": &types.AttributeValueMemberS{Value: testutil.ValidUser1.Email},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualAttribute := GetKey(tt.user)

			assert.Equal(t, tt.expectedAttribute, actualAttribute)
		})
	}
}
