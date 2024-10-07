# Question Service

## Overview

The Question Service is built with Go, utilizing Firestore as the database and Chi as the HTTP router. It allows for basic operations such as creating, reading, updating, and deleting question records.

## Features

- Create new questions
- Retrieve question information by ID
- Update question details
- Delete question

## Technologies Used

- Go (Golang)
- Firestore (Google Cloud Firestore)
- Chi (HTTP router)

## Getting Started

### Prerequisites

- Go 1.16 or later
- Google Cloud SDK
- Firestore database setup in your Google Cloud project

### Installation

1. Clone the repository

2. Set up your Firestore client

3. Install dependencies:

```bash
go mod tidy
```

4. Create the `.env` file from copying the `.env.example`, and fill in the `FIREBASE_CREDENTIAL_PATH` with the path of the firebase credential JSON file.

### Running the Application

To start the server, run the following command:

```bash
go run main.go
```

The server will be available at http://localhost:8080.

## Running the Application via Docker

To run the application via Docker, run the following command:

```bash
docker build -t question-service .
```

```bash
docker run -p 8080:8080 --env-file .env -d question-service
```

The server will be available at http://localhost:8080.

## API Endpoints

- `POST /questions`
- `GET /questions/{id}`
- `GET /questions`
- `PUT /questions/{id}`
- `DELETE /questions/{id}`

## Managing Firebase

To reset and repopulate the database, run the following command:

```bash
go run main.go -populate
```
