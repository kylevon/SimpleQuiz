# Quiz Application

This project consists of a backend API server and a CLI client for a quiz application.

## Prerequisites

- Go 1.23.1 or later

## Project Structure

```
.
├── backend/
│   ├── main.go
│   ├── main_test.go
│   ├── quiz.json
│   └── go.mod
├── cli/
│   ├── cmd/
│   │   ├── get.go
│   │   ├── post.go
│   │   └── root.go
│   ├── main.go
│   ├── go.mod
│   └── go.sum
└── README.md
```

## Running the Backend

1. Navigate to the backend directory:
   ```
   cd backend
   ```

2. Ensure that the `quiz.json` file is present in the backend directory. This file contains the quiz questions, answer options, and correct answers.

3. Run the server:
   ```
   go run main.go
   ```

The server will start on `http://localhost:8080`. It will load the quiz data from `quiz.json` on startup.

## Running the CLI

1. Navigate to the cli directory:
   ```
   cd cli
   ```

2. Build the CLI:
   ```
   go build -o quiz-cli
   ```

3. Run the CLI commands:
   - To get quiz questions:
     ```
     ./quiz-cli get
     ```
   - To submit answers:
     ```
     ./quiz-cli post 2 1 1 1 2
     ```

## Running Tests

To run the tests for the backend:

1. Navigate to the backend directory:
   ```
   cd backend
   ```

2. Run the tests:
   ```
   go test
   ```

The `main_test.go` file contains unit tests for the backend. It tests the following functions:

- `TestParseAndValidateAnswers`: Tests the input validation for quiz submissions.
- `TestPostQuiz`: Tests the quiz submission endpoint with various inputs.

The test file also sets up a temporary test quiz file and cleans it up after the tests are complete.

## Note

- Make sure the backend server is running when using the CLI to interact with the quiz API.
- The backend relies on the `quiz.json` file for quiz data. Ensure this file is present and correctly formatted in the backend directory before starting the server.