# Execution Service

### Installation

1. Install dependencies:

```bash
go mod tidy
```

2. Create the `.env` file from copying the `.env.example`, and copy the firebase JSON file into execution-service/ fill in the `FIREBASE_CREDENTIAL_PATH` with the path of the firebase credential JSON file.

### Running the Application

To start the server, run the following command:

```bash
go run main.go
```

The server will be available at http://localhost:8083.

### Setting up message queue with RabbitMQ

A message queue is used to pass submission results asynchronously from the execution-service to the history-service.

1. In order to do so, we can run the following command to set up a docker container for RabbitMQ:

```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

2. Then we can run the execution-service:

```bash
go run main.go
```

3. We can run the history-service by changing our directory and running the same command:

```bash
cd ../history-service
go run main.go
```

To view more details on the RabbitMQ queue, we can go to `localhost:15672`, and login using `guest` as the username and password.

### Running the Application via Docker

To run the application via Docker, run the following command:

```bash
docker build -t execution-service .
```

```bash
docker run -p 8083:8083 --env-file .env -d execution-service
```

The server will be available at http://localhost:8083.

## API Endpoints

- `POST /tests/populate`
- `GET /tests/{questionDocRefId}/`
- `POST /tests/{questionDocRefId}/execute`
- `POST /tests/{questionDocRefId}/submit`

## Managing Firebase

To reset and repopulate the database, run the following command:

```bash
go run main.go
```

## Repopulate test cases

To repopulate test cases, you need to repopulate the questions in the question-service, which will automatically call the execution-service to populate the test cases.

In question-service, run the following command:

```bash
go run main.go -populate
```

## API Documentation

`GET /tests/{questionDocRefId}/`

To read visible test cases via a question ID, run the following command:

```bash
curl -X GET http://localhost:8083/tests/{questioinDocRefId}/ \
-H "Content-Type: application/json"
```

The following json format will be returned:

```json
[
  {
    "input": "hello",
    "expected": "olleh"
  }
]
```

`POST /tests`

To create a new test case, run the following command:

```bash
curl -X POST http://localhost:8083/tests \
-H "Content-Type: application/json" \
-d '{
"questionDocRefId": "sampleDocRefId123",
"questionTitle": "Sample Question Title",
"visibleTestCases": "2\nhello\nolleh\nHannah\nhannaH",
"hiddenTestCases": "2\nHannah\nhannaH\nabcdefg\ngfedcba"
}'
```

`PUT /tests/{questionDocRefId}`

To update an existing test case from an existing question, run the following command:

```bash
curl -X PUT http://localhost:8083/tests/{questionDocRefId} \
-H "Content-Type: application/json" \
-d '{
"visibleTestCases": "2\nhello\nolleh\nHannah\nhannaH",
"hiddenTestCases": "2\nHannah\nhannaH\nabcdefg\ngfedcba"
}'
```

`DELETE /tests/{questionDocRefId}`

To delete an existing test case from an existing question, run the following command:

```bash
curl -X DELETE http://localhost:8083/tests/{questionDocRefId} \
-H "Content-Type: application/json"
```

`POST /tests/{questionDocRefId}/execute`

To execute test cases via a question ID without custom test cases, run the following command, with custom code and language:

```bash
curl -X POST http://localhost:8083/tests/{questioinDocRefId}/execute \
-H "Content-Type: application/json" \
-d '{
"code": "name = input()\nprint(name[::-1])",
"language": "Python"
}'
```

The following json format will be returned:

```json
{
  "visibleTestResults": [
    {
      "input": "hello",
      "expected": "olleh",
      "actual": "olleh",
      "passed": true,
      "error": ""
    }
  ],
  "customTestResults": null
}
```

To execute visible and custom test cases via a question ID with custom test cases, run the following command, with custom code, language and custom test cases:

```bash
curl -X POST http://localhost:8083/tests/{questioinDocRefId}/execute \
-H "Content-Type: application/json" \
-d '{
"code": "name = input()\nprint(name[::-1])",
"language": "Python",
"customTestCases": "2\nHannah\nhannaH\nabcdefg\ngfedcba\n"
}'
```

The following json format will be returned:

```json
{
  "visibleTestResults": [
    {
      "input": "hello",
      "expected": "olleh",
      "actual": "olleh",
      "passed": true,
      "error": ""
    }
  ],
  "customTestResults": [
    {
      "input": "Hannah",
      "expected": "hannaH",
      "actual": "hannaH",
      "passed": true,
      "error": ""
    },
    {
      "input": "abcdefg",
      "expected": "gfedcba",
      "actual": "gfedcba",
      "passed": true,
      "error": ""
    }
  ]
}
```

`POST /tests/{questionDocRefId}/submit`

To submit a solution and execute visible and hidden test cases via a question ID, run the following command, with custom code and language:

```bash
curl -X POST http://localhost:8083/tests/{questioinDocRefId}/submit \
-H "Content-Type: application/json" \
-d '{
"title": "Example Title",
"code": "name = input()\nprint(name[::-1])",
"language": "Python",
"user": "user123",
"matchedUser": "user456",
"matchedTopics": ["topic1", "topic2"],
"questionDifficulty": "Medium",
"questionTopics": ["Loops", "Strings"]
}'
```

The following json format will be returned:

```json
{
  "visibleTestResults": [
    {
      "input": "hello",
      "expected": "olleh",
      "actual": "olleh",
      "passed": true,
      "error": ""
    }
  ],
  "hiddenTestResults": {
    "passed": 2,
    "total": 2
  },
  "status": "Accepted"
}
```

If compilation error exists or any of the tests (visible and hidden) fails, status "Attempted" will be returned:

```json
{
  "visibleTestResults": [
    {
      "input": "hello",
      "expected": "olleh",
      "actual": "",
      "passed": false,
      "error": "Command execution failed: Traceback (most recent call last):\n  File \"/tmp/4149249165.py\", line 2, in \u003cmodule\u003e\n    prit(name[::-1])\n    ^^^^\nNameError: name 'prit' is not defined. Did you mean: 'print'?\n: %!w(*exec.ExitError=\u0026{0x4000364678 []})"
    }
  ],
  "hiddenTestResults": {
    "passed": 0,
    "total": 2
  },
  "status": "Attempted"
}
```
