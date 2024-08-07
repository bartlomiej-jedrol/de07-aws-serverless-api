START_TIME ?=
AWS_LAMBDA_FUNCTION = "de07-lambda"
AWS_LAMBDA_LOG_GROUP = "/aws/lambda/de07-lambda"
AWS_LAMBDA_BUILD_ZIP_PATH = "build/main.zip"

clean:
	rm -f bootstrap $(AWS_LAMBDA_BUILD_ZIP_PATH)

build: clean
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap cmd/main.go
	zip $(AWS_LAMBDA_BUILD_ZIP_PATH) bootstrap

push: build
	aws lambda update-function-code --function-name $(AWS_LAMBDA_FUNCTION) \
	--zip-file fileb://$(AWS_LAMBDA_BUILD_ZIP_PATH)

logs:
	echo $(START_TIME)
	@sleep 2
	aws logs filter-log-events \
	--log-group-name $(AWS_LAMBDA_LOG_GROUP) \
	--start-time $(START_TIME) \
	--limit 50 \
	--color auto \
	--output text

test:
	curl --header "Content-Type: application/json" --request POST --data '{"email": "bartlomiej.jedrol@gmail.com", "firstName": "Bartlomiej", "lastName": "Jedrol", "age": 37}' https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev