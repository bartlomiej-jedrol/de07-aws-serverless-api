// Main implements an entry point of the Lambda function.
package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/handlers"
)

func HandleRequest(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("========== Request ==========: %v", request)
	log.Printf("========== HTTPMethod ==========: %v", request.HTTPMethod)
	log.Printf("========== Headers ==========: %v", request.Headers)
	log.Printf("========== PathParameters ==========: %v", request.PathParameters)
	log.Printf("========== QueryStringParameters ==========: %v", request.QueryStringParameters)
	log.Printf("========== Body ==========: %v", request.Body)
	switch request.HTTPMethod {
	case "GET":
		return handlers.GetUser(request)
	case "POST":
		return handlers.CreateUser(request)
	// case "PUT":
	// 	return handlers.UpdateUser(request, tableName)
	// case "DELETE":
	// 	return handlers.DeleteUser(request, tableName)
	default:
		return handlers.UnhandledHTTPMethod(request)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
