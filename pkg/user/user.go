// User implements functions for interacting with DynamoDB database.
package user

import (
	"context"
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
)

func init() {
	// Load AWS config (~/.aws/config).
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Create DynamoDB client.
	userTable.DynamoDbClient = dynamodb.NewFromConfig(cfg)
	if userTable.DynamoDbClient == nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}

	// Create validator of User struct.
	validate = validator.New()
}

func FetchUser(email string) (*User, error) {
	// Build input with key (user's email) of the item to be fetched from the DynamoDB table.
	key := map[string]types.AttributeValue{"email": &types.AttributeValueMemberS{Value: email}}
	input := &dynamodb.GetItemInput{Key: key, TableName: &userTable.TableName}

	// Get user data from the DynamoDB table. If err return to the caller.
	response, err := userTable.DynamoDbClient.GetItem(context.TODO(), input)
	log.Printf("========== response ==========: %v", response)
	if err != nil {
		log.Printf("Failed to get item (user data) from the DynamoDB API. Error: %v", err)
		return nil, err
	}

	// Extract user data from the DynamoDB output.
	var user User
	err = attributevalue.UnmarshalMap(response.Item, &user)
	if err != nil {
		log.Printf("Failed to unmarshal map for item get from the DynamoDB table: %v. Error: %v", response.Item, err)
	}
	log.Printf("========== user ==========: %v", user)

	// Validate the user struct.
	err = validate.Struct(user)
	if err != nil {
		log.Printf("Failed to validate the user: %v. Error: %v", user, err)
		return nil, err
	}
	return &user, nil
}

func FetchUsers() {

}

// CreateUser creates user in DynamoDB table. I returns error in case of failure.
func CreateUser(user User) error {
	// Prepare user item with all attributes.
	item := map[string]types.AttributeValue{
		"email":     &types.AttributeValueMemberS{Value: user.Email},
		"firstName": &types.AttributeValueMemberS{Value: user.FirstName},
		"lastName":  &types.AttributeValueMemberS{Value: user.LastName},
		"age":       &types.AttributeValueMemberN{Value: strconv.Itoa(user.Age)},
	}
	log.Printf("========== item ==========: %v", item)

	input := dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(userTable.TableName)}
	log.Printf("========== input ==========: %v", input)

	// Put item into DynamoDB table.
	_, err := userTable.DynamoDbClient.PutItem(context.TODO(), &input)
	if err != nil {
		log.Printf("Failed to put item to the DynamoDb table: %v", err)
		return err
	}
	return nil
}

func UpdateUser() {
}

func DeleteUser() {

}
