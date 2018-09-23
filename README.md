# slim-load-recorder

## Overview
Work In Progress

## Requirements
* Serverless Framework
    
    `npm install -g serverless`
* AWS developer account

## Usage
1. Build
`$ env GOOS=linux go build -o bin/main`

2. Deploy
`$ serverless deploy`

3. After finising upload your app, access to your endpoint 

### Run in local (Optional)
1. Install

`$ npm install -g aws-sam-local`

2. Run in local

`$ env APP_ENV=local sam local start-api`

3. Access to your local endpoint

`$ curl http://localhost:3000/slim`
