# slim-load-recorder
[![Build Status](https://travis-ci.org/kutsuzawa/slim-load-recorder.svg?branch=master)](https://travis-ci.org/kutsuzawa/slim-load-recorder) [![Maintainability](https://api.codeclimate.com/v1/badges/33d4eb3e51099945979d/maintainability)](https://codeclimate.com/github/kutsuzawa/slim-load-recorder/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/33d4eb3e51099945979d/test_coverage)](https://codeclimate.com/github/kutsuzawa/slim-load-recorder/test_coverage)  
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
`$ curl -X POST <YOUR ENDPOINT>:3000/loads -H 'Content-Type:application/json' -d '{"user_id": "bookun", "weight":71.8, "distance":5.2, "date": "2018-09-24 07:00:00", "start_at": "2018-09-23 00:00:00", "end_at":"2018-09-30 00:00:00"}'
`

### Run in local (Optional)
1. Install

`$ npm install -g aws-sam-local`

2. Run in local
```
$ make setup
$ make local
```

3. Access to your local endpoint

`$ curl -X POST http://localhost:3000/loads -H 'Content-Type:application/json' -d '{"user_id": "bookun", "weight":71.8, "distance":5.2, "date": "2018-09-24 07:00:00", "start_at": "2018-09-23 00:00:00", "end_at":"2018-09-30 00:00:00"}'`
