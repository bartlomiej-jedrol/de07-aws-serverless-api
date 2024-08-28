// User implements functions for interacting with DynamoDB database.
package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	validator "github.com/go-playground/validator/v10"
)

type User struct {
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

var (
	userTable = TableBasics{TableName: "de07-user"}
	validate  *validator.Validate

	// ErrorFailedToLoadAWSConfig Failed to load AWS config.
	ErrorFailedToLoadAWSConfig = errors.New("failed to load AWS config")

	// FailedToCreateDynamoDBClient Failed to create DynamoDB client.
	ErrorFailedToCreateDynamoDBClient = errors.New("failed to create DynamoDB client")

	// ErrorFailedToUnmarshalMap Failed to unmarshal map.
	ErrorFailedToUnmarshalMap = errors.New("failed to unmarshal map for item")

	// ErrorFailedToValidateUser Failed to validate user.
	ErrorFailedToValidateUser = errors.New("failed to validate user")

	// ErrorUserDoesNotExist User does not exist.
	ErrorUserDoesNotExist = errors.New("user does not exist")

	//ErrorFailedToGetItem Failed to get item.
	ErrorFailedToGetItem = errors.New("failed to get item from DynamoDB")

	// ErrorFailedToPutItem Failed to put item.
	ErrorFailedToPutItem = errors.New("failed to put item to DynamoDB")

	// ErrorFailedToDeleteItem Failed to delete item.
	ErrorFailedToDeleteItem = errors.New("failed to delete item from DynamoDB")
)

func init() {
	// Load AWS config (~/.aws/config).
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("%v: %v", ErrorFailedToLoadAWSConfig, err)
	}

	// Create DynamoDB client.
	userTable.DynamoDbClient = dynamodb.NewFromConfig(cfg)
	if userTable.DynamoDbClient == nil {
		log.Fatalf("%v: %v", ErrorFailedToCreateDynamoDBClient, err)
	}

	// Create validator of User struct.
	validate = validator.New()
}

// FetchUser fetches provided item from DynamoDB table based on key (email).
func FetchUser(email string) (*User, error) {
	// Build input with key (user's email).
	input := dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: email}},
		TableName: aws.String(userTable.TableName),
	}

	// Get user data from DynamoDB table.
	r, err := userTable.DynamoDbClient.GetItem(context.TODO(), &input)
	log.Printf("FetchUser response: %v", r)
	if err != nil {
		log.Printf("%v: %v", ErrorFailedToGetItem, err)
		return nil, ErrorFailedToGetItem
	}

	// Return an error if user does not exist (r.Item is nil).
	if r.Item == nil {
		log.Printf("%v: %v", ErrorUserDoesNotExist, err)
		return nil, ErrorUserDoesNotExist
	}

	// Extract user data from DynamoDB output.
	var u User
	err = attributevalue.UnmarshalMap(r.Item, &u)
	if err != nil {
		log.Printf("%v: %v, %v", ErrorFailedToUnmarshalMap, r.Item, err)
		return nil, ErrorFailedToUnmarshalMap
	}
	log.Printf("user: %v", u)

	return &u, nil
}

// FetchUsers fetches items from DynamoDB table.
func FetchUsers() ([]User, error) {
	var input *dynamodb.BatchGetItemInput
	var users []User
	r, err := userTable.DynamoDbClient.BatchGetItem(context.TODO(), input)
	fmt.Printf("%v, %v", r, err)
	return users, nil
}

// CreateUser creates user in DynamoDB table.
func CreateUser(user User) error {
	// Prepare user item with all attributes.
	item := map[string]types.AttributeValue{
		"email":     &types.AttributeValueMemberS{Value: user.Email},
		"firstName": &types.AttributeValueMemberS{Value: user.FirstName},
		"lastName":  &types.AttributeValueMemberS{Value: user.LastName},
		"age":       &types.AttributeValueMemberN{Value: strconv.Itoa(user.Age)},
	}
	log.Printf("CreateUser item: %v", item)

	// Prepare input for PutItem method.
	input := dynamodb.PutItemInput{
		Item:         item,
		TableName:    aws.String(userTable.TableName),
		ReturnValues: "ALL_OLD",
	}
	log.Printf("CreateUser input: %v", input)

	// Put item into DynamoDB table.
	_, err := userTable.DynamoDbClient.PutItem(context.TODO(), &input)
	if err != nil {
		log.Printf("%v: %v", ErrorFailedToPutItem, err)
		return ErrorFailedToPutItem
	}

	// Logging methods.
	// Unmarshaling an entire map.
	// var responseUser User
	// err = attributevalue.UnmarshalMap(r.Attributes, &responseUser)
	// if err != nil {
	// 	log.Printf("%v: %v", ErrorFailedToUnmarshalMap, err)
	// }
	// log.Printf("responseAttributes: %v", responseUser)

	// Unmarshaling a single attribute.
	// var userEmail string
	// err = attributevalue.Unmarshal(item["email"], &userEmail)
	// if err != nil {
	// 	log.Printf("%v: %v", ErrorFailedToUnmarshalMap, err)
	// }
	// log.Printf("unmarshal: %v", userEmail)

	// // Printing value of a single attribute.
	// log.Printf("item: %v", item["email"].(*types.AttributeValueMemberS).Value)

	return nil
}

// UpdateUser updates existing user in DynamoDB table.
func UpdateUser(user User) error {
	// Validate user struct if it has required email field.
	err := validate.Struct(user)
	if err != nil {
		log.Printf("%v: %v, %v", ErrorFailedToValidateUser, user, err)
		return ErrorFailedToValidateUser
	}

	var u *User
	u, err = FetchUser(user.Email)
	if err != nil {
		return err // Bypassing error from the FetchUser function to the caller to build response.
	}

	// If the user exist create it again to overwrite data.
	if u != nil {
		err := CreateUser(user)
		if err != nil {
			return err // Bypassing error from the FetchUser function to the caller to build response.
		}
	}

	return nil
}

// DeleteUser deletes provided item to be deleted from DynamoDB table based on key (email).
func DeleteUser(email string) (*User, error) {
	// Check for user existence.
	u, err := FetchUser(email)
	if err != nil {
		return nil, err // Bypassing error from the FetchUser function to the caller to build response.
	}

	// Build input with key (user's email).
	input := dynamodb.DeleteItemInput{
		Key:          map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: email}},
		TableName:    aws.String(userTable.TableName),
		ReturnValues: "ALL_OLD",
	}
	log.Printf("DeleteUser input: %v", input)

	// Delete item from DynamoDB table.
	r, err := userTable.DynamoDbClient.DeleteItem(context.TODO(), &input)
	log.Printf("DeleteItem response, err: %v: %v", r, err)
	if err != nil {
		log.Printf("%v: %v", ErrorFailedToDeleteItem, err)
		return nil, ErrorFailedToDeleteItem
	}

	// Return an error if user does not exist (r.Attributes is nil).
	if r.Attributes == nil {
		log.Printf("user does not exist: %v", email)
		return nil, ErrorUserDoesNotExist
	}

	return u, nil
}

// GetKey returns key of a user in a required format.
func GetKey(user User) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: user.Email}}
}
