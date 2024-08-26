// User implements functions for interacting with DynamoDB database.
package user

import (
	"context"
	"errors"
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

	//ErrorFailedToGetItem Failed to get item.
	ErrorFailedToGetItem = errors.New("failed to get item (user data) from the DynamoDB")

	// ErrorFailedToUnmarshalMap Failed to unmarshal map.
	ErrorFailedToUnmarshalMap = errors.New("failed to unmarshal map for item get from the DynamoDB, item")

	// ErrorFailedToValidateUser Failed to validate user.
	ErrorFailedToValidateUser = errors.New("failed to validate the user")

	// ErrorFailedToPutItem Failed to put item.
	ErrorFailedToPutItem = errors.New("failed to put item to the DynamoDB")
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

func FetchUser(email string) (*User, error) {
	// Build input with key (user's email) of the item to be fetched from the DynamoDB table.
	key := map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: email}}
	input := dynamodb.GetItemInput{Key: key, TableName: aws.String(userTable.TableName)}

	// Get user data from the DynamoDB table. If err return to the caller.
	response, err := userTable.DynamoDbClient.GetItem(context.TODO(), &input)
	log.Printf("========== response ==========: %v", response)
	if err != nil {
		log.Printf("%v: %v", ErrorFailedToGetItem, err)
		return nil, err
	}

	// Extract user data from the DynamoDB output.
	var user User
	err = attributevalue.UnmarshalMap(response.Item, &user)
	if err != nil {
		log.Printf("%v: %v, %v", ErrorFailedToUnmarshalMap, response.Item, err)
		return nil, err
	}
	log.Printf("========== user ==========: %v", user)

	// Validate the user struct.
	err = validate.Struct(user)
	if err != nil {
		log.Printf("%v: %v, %v", ErrorFailedToValidateUser, user, err)
		return nil, err
	}
	return &user, nil
}

func FetchUsers(user User) {
}

// CreateUser creates user in DynamoDB table. It returns error in case of failure.
func CreateUser(user User) error {
	// Prepare user item with all attributes.
	item := map[string]types.AttributeValue{
		"email":     &types.AttributeValueMemberS{Value: user.Email},
		"firstName": &types.AttributeValueMemberS{Value: user.FirstName},
		"lastName":  &types.AttributeValueMemberS{Value: user.LastName},
		"age":       &types.AttributeValueMemberN{Value: strconv.Itoa(user.Age)},
	}
	log.Printf("========== item ==========: %v", item)

	// Prepare input for the PutItem method.
	input := dynamodb.PutItemInput{
		Item:         item,
		TableName:    aws.String(userTable.TableName),
		ReturnValues: "ALL_OLD",
	}
	log.Printf("========== input ==========: %v", input)

	// Put item into DynamoDB table.
	_, err := userTable.DynamoDbClient.PutItem(context.TODO(), &input)
	if err != nil {
		log.Printf("%v: %v", ErrorFailedToPutItem, err)
		return err
	}

	// Logging methods.
	// Unmarshaling an entire map.
	// var responseUser User
	// err = attributevalue.UnmarshalMap(r.Attributes, &responseUser)
	// if err != nil {
	// 	log.Printf("%v: %v", ErrorFailedToUnmarshalMap, err)
	// }
	// log.Printf("========== responseAttributes ==========: %v", responseUser)

	// Unmarshaling a single attribute.
	// var userEmail string
	// err = attributevalue.Unmarshal(item["email"], &userEmail)
	// if err != nil {
	// 	log.Printf("%v: %v", ErrorFailedToUnmarshalMap, err)
	// }
	// log.Printf("========== unmarshal ==========: %v", userEmail)

	// // Printing value of a single attribute.
	// log.Printf("========== item ==========: %v", item["email"].(*types.AttributeValueMemberS).Value)

	return nil
}

// UpdateUser updates existing user in DynamoDB table. It returns error in case of failure.
func UpdateUser(user User) error {
	u, _ := FetchUser(user.Email)
	if u != nil {
		err := CreateUser(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteUser() {

}

// GetKey returns key of a user in a required format.
func GetKey(user User) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: user.Email}}
}
