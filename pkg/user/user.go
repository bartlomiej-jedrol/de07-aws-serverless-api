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
	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/models"
	validator "github.com/go-playground/validator/v10"
)

var (
	userTable = models.TableBasics{TableName: "de07-user"}
	validate  *validator.Validate

	ErrorFailedToLoadAWSConfig        = errors.New("failed to load AWS config")
	ErrorFailedToCreateDynamoDBClient = errors.New("failed to create DynamoDB client")
	ErrorFailedToUnmarshalMap         = errors.New("failed to unmarshal map for item")
	ErrorFailedToValidateUser         = errors.New("failed to validate user")
	ErrorUserDoesNotExist             = errors.New("user does not exist")
	ErrorFailedToGetItem              = errors.New("failed to get item from DynamoDB")
	ErrorFailedToGetItems             = errors.New("failed to get items from DynamoDB")
	ErrorFailedToPutItem              = errors.New("failed to put item to DynamoDB")
	ErrorFailedToDeleteItem           = errors.New("failed to delete item from DynamoDB")
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
func FetchUser(email string) (*models.User, error) {
	// Build input with key (user's email).
	input := dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"email": &types.AttributeValueMemberS{Value: email},
		},
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
	var u models.User
	err = attributevalue.UnmarshalMap(r.Item, &u)
	if err != nil {
		log.Printf("%v: %v, %v", ErrorFailedToUnmarshalMap, r.Item, err)
		return nil, ErrorFailedToUnmarshalMap
	}
	log.Printf("user: %v", u)

	return &u, nil
}

// FetchUsers fetches items from DynamoDB table.
func FetchUsers() ([]models.User, error) {
	// Scan items of DynamoDB table.
	input := dynamodb.ScanInput{TableName: &userTable.TableName}
	r, err := userTable.DynamoDbClient.Scan(context.TODO(), &input)
	if err != nil {
		log.Printf("%v: %v", ErrorFailedToGetItems, err)
		return nil, ErrorFailedToGetItems
	}
	log.Printf("r.Items: %v", r.Items)

	// Build list of users.
	var users []models.User
	for _, item := range r.Items {
		var u models.User
		err := attributevalue.UnmarshalMap(item, &u)
		if err != nil {
			log.Printf("%v: %v", ErrorFailedToUnmarshalMap, err)
			return nil, ErrorFailedToUnmarshalMap
		}
		users = append(users, u)
	}
	log.Printf("users: %v", users)

	return users, nil
}

// CreateUser creates user in DynamoDB table.
// It does not return created user - instead the user is taken from the API body request.
func CreateUser(user models.User) error {
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
// It does not return updated user - instead the user is taken from the API body request.
func UpdateUser(user models.User) error {
	// Validate user struct if it has required email field.
	err := validate.Struct(user)
	if err != nil {
		log.Printf("%v: %v, %v", ErrorFailedToValidateUser, user, err)
		return ErrorFailedToValidateUser
	}

	var u *models.User
	u, err = FetchUser(user.Email)
	if err != nil {
		return err // Bypassing error from the FetchUser function to the caller to build response.
	}

	// If the user exist create it again to overwrite data.
	if u != nil {
		err := CreateUser(user)
		if err != nil {
			return err // Bypassing error from FetchUser function to the caller to build response.
		}
	}

	return nil
}

// DeleteUser deletes provided item to be deleted from DynamoDB table based on key (email).
func DeleteUser(email string) (*models.User, error) {
	// Check for user existence.
	u, err := FetchUser(email)
	if err != nil {
		return nil, err // Bypassing error from FetchUser function to the caller to build response.
	}

	// Build input with key (user's email).
	input := dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"email": &types.AttributeValueMemberS{Value: email},
		},
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
func GetKey(user models.User) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"email": &types.AttributeValueMemberS{Value: user.Email},
	}
}
