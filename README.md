# slim-load-recorder

## Overview
Work In Progress

## Requirements
* Serverless Framework

`$ npm install -g serverless`
* AWS developer account

## Usage

```
$ make setup
$ make lambda
```
After finising upload your app, post to your endpoint.
ex)
`$ curl -X POST <YOUR ENDPOINT>:3000/add -H 'Content-Type:application/json' -d '{"weight":71.8, "distance":5.2}'`

### Run in local (Optional)
1. Install

`$ npm install -g aws-sam-local`

2. Run in local
```
$ make setup
$ make local
```

3. Access to your local endpoint

`$ curl -X POST http://localhost:3000/add -H 'Content-Type:application/json' -d '{"weight":71.8, "distance":5.2}' `
