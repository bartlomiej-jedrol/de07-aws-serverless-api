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
	@sleep 5
	aws logs filter-log-events \
	--log-group-name $(AWS_LAMBDA_LOG_GROUP) \
	--start-time $(START_TIME) \
	--limit 10000 \
	--color auto \
	--output text


# GET

get:
	curl --header "Content-Type: application/json" --request GET https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev\?email\=bartlomiej.jedrol@gmail.com

get_user_does_not_exist:
	curl --header "Content-Type: application/json" --request GET https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev\?email\=test.test@gmail.com

get_empty_query_parameter:
	curl --header "Content-Type: application/json" --request GET https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev\?email\=

get_all:
	curl --header "Content-Type: application/json" --request GET https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev

# POST

post:
	curl --header "Content-Type: application/json" --request POST --data '{"email": "bartlomiej.jedrol@gmail.com", "firstName": "Bartlomiej", "lastName": "Jedrol", "age": 37}' https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev

post_invalid_json:
	curl --header "Content-Type: application/json" --request POST --data '{""email": "jedrol.natalia@gmail.com", "firstName": "Natalia", "lastName": "Jedrol", "age": 33}'  https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev

post_missing_email:
	curl --header "Content-Type: application/json" --request POST --data '{"firstName": "Natalia", "lastName": "Jedrol", "age": 33}' https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev

# PUT

put:
	curl --header "Content-Type: application/json" --request PUT --data '{"email": "bartlomiej.jedrol@gmail.com", "firstName": "BartlomiejUpdated4", "lastName": "JedrolUpdated4", "age": 37}' https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev\?email\=bartlomiej.jedrol@gmail.com

put_invalid_json:
	curl --header "Content-Type: application/json" --request PUT --data '{""email": "bartlomiej.jedrol@gmail.com", "firstName": "BartlomiejUpdated4", "lastName": "JedrolUpdated4", "age": 37}' https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev\?email\=bartlomiej.jedrol@gmail.com

put_missing_email:
	curl --header "Content-Type: application/json" --request PUT --data '{"firstName": "BartlomiejUpdated4", "lastName": "JedrolUpdated4", "age": 37}' https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev\?email\=bartlomiej.jedrol@gmail.com

put_user_does_not_exist:
	curl --header "Content-Type: application/json" --request PUT --data '{"email": "test.test@gmail.com", "firstName": "test", "lastName": "test", "age": 1}' https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev\?email\=test.test@gmail.com

# DELETE

delete:
	curl --header "Content-Type: application/json" --request DELETE https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev\?email\=bartlomiej.jedrol@gmail.com

delete_empty_query_parameter:
	curl --header "Content-Type: application/json" --request GET https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev\?email\=

delete_user_does_not_exist:
	curl --header "Content-Type: application/json" --request DELETE https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev\?email\=test.test@gmail.com

# OTHER

not_allowed:
	curl --header "Content-Type: application/json" --request PATCH --data '{"email": "bartlomiej.jedrol@gmail.com", "firstName": "Bartlomiej_updated", "lastName": "Jedrol_updated", "age": 37_updated}' https://7t5wi1q5p4.execute-api.eu-central-1.amazonaws.com/dev\?email\=bartlomiej.jedrol@gmail.com