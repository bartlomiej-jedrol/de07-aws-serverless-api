// User implements functions for interacting with DynamoDB database.
package user

import (
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

var userTable = TableBasics{TableName: "de07-user"}

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
}

func FetchUser() {

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
