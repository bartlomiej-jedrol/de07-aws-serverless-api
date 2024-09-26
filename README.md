# Architecture

![Architecture](architecture.png)

# Project Description

When stepping into Serverless 🌐, finding the right architecture 🏛️ can make all the difference. 



I recently built a Serverless API on AWS 🚀implementing the Services Design Pattern. In the Services Pattern, a single Lambda function can handle a few jobs that are usually related to a single entity of the data model (e.g., User 👤). 


All CRUD operations on the User data model are performed on the single HTTP endpoint using different HTTP methods. 


For this to work, you can have a 'small' router (HTTP method) at the beginning of your Lambda code.⚙️ 


It's different from the Microservices Pattern where each job is isolated within a separate Lambda function.💡 



Benefits of the Services Design Pattern: ✔️
- Fewer Lambda functions to manage ✔️
- Reduced cold starts ✔️
- Team autonomy 👥
- Faster deployments 🚀



Drawbacks of Services Pattern: ❌
- More complicated debugging ❌
-  Requires a router ❌
-  Bigger function size ❌



Tech Stack I utilised: 💻
- Go
- Terraform



AWS Services I utilised: 🛠️
- AWS API Gateway 
- AWS Lambda 
- AWS DynamoDB 
- AWS CloudWatch