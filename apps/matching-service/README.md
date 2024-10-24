# Matching Service

The Matching Service provides a WebSocket server to manage real-time matching requests through WebSocket connections. It is implemented in Go and uses the Gorilla WebSocket package to facilitate communication.

## Setup

### Prerequisites

Ensure you have Go installed on your machine.

### Installation

1. Navigate to the matching service directory:

```bash
cd ./apps/matching-service
```

2. Install the necessary dependencies:

```bash
go mod tidy
```

3. Create a copy of the `.env.example` file as an `.env` file with the following environment variables:

- `PORT`: Specifies the port for the WebSocket server. Default is `8081`.
- `JWT_SECRET`: The secret key used to verify JWT tokens.
- `MATCH_TIMEOUT`: The time in seconds to wait for a match before timing out.
- `REDIS_URL`: The URL for the Redis server. Default is `localhost:6379`. If you are using docker, use `redis-container:6379`

4. Start a local Redis server:

```bash
docker run -d -p 6379:6379 redis
```

5. Start the WebSocket server:

```bash
go run main.go
```

## API Usage

To establish a WebSocket connection with the matching service, use the following JavaScript code:

```javascript
const ws = new WebSocket("ws://localhost:8081/match");
```

### Authentication

The initial WebSocket request should include a JWT token that contains the user’s ID as a claim. This token will be verified by the server to authenticate the user. The user ID extracted is used to identify the client and facilitate the matching process.

### Matching Workflow

1. **Sending Matching Parameter**s: Once the WebSocket connection is established, the client sends a message containing the matching parameters (e.g., preferred topics or difficulty levels).

2. **Receiving Match Results**:
   2.1. **Match Found**: If a match is found, the server sends the matching result back to the client via the WebSocket connection.
   2.2. **No Match Found**: If after a set amount of time, no match is found, the request timeouts, and the server sends a message that the matching failed.

3. **Connection Closure**:
   3.1. **Received Match Result**: After sending the match result, the server closes the connection.
   3.2. **Cancellation**: If the user cancels the matching process, the client should close the connection. The server will recognize this and cancel the ongoing match.

### Message Formats

Provide the structure of the messages being sent back and forth between the server and the client. This includes the format of the initial matching parameters and the response payload. All requests should be in JSON and contain the `type` field to handle different types of messages.

Client sends matching parameters:

```json
{
  "type": "match_request",
  "topics": ["Algorithms", "Arrays"],
  "difficulties": ["Easy", "Medium"],
  "username": "Jane Doe"
}
```

Server response on successful match:

```json
{
  "type": "match_found",
  "matchId": "1c018916a34c5bee21af0b2670bd6156",
  "user": "zkb4px",
  "matchedUser": "JohnDoe",
  "topic": "Algorithms",
  "difficulty": "Medium"
}
```

If no match is found after a set period of time, the server will send a timeout message:

```json
{
  "type": "timeout",
  "message": "No match found. Please try again later."
}
```

If user has an existing websocket connection and wants to initiate another match, the server will reject the request:

```json
{
  "type": "match_rejected",
  "message": "You are already in a matchmaking queue. Please disconnect before reconnecting."
}
```

If the server encounters an issue during the WebSocket connection or processing, the connection will be closed without any error message. The client should treat the unexpected closing as an error.

## Testing

Utilize `./tests/websocket-test.html` for a basic debugging interface of the matching service. This HTML file provides an interactive way to test the WebSocket connection, send matching requests, and observe responses from the server.

Make sure to open the HTML file in a web browser while the WebSocket server is running to perform your tests.

You can open one instance of the HTML file in multiple tabs to simulate multiple clients connecting to the server. (In the future: ensure that only one connection is allowed per user)

## Running the Application via Docker

Before running the following commands, ensure that the URL for the Redis server in `.env` file has been changed to `REDIS_URL=redis-container:6379`

To run the application via Docker, run the following command:

1. Set up the Go Docker container for the matching service
```bash
docker build -f Dockerfile -t matching-service .
```

2. Create the Docker network for matching service
```bash
docker network create matching-service-network
```

3. Start a new Redis container in detached mode using the Redis image from Docker Hub, and run it in the matching service network
```bash
docker run -d --name redis-container --network matching-service-network redis
```

4. Run the Go Docker container for the matching-service inside the same matching service network as the redis container, so that they can communicate with one another
```bash
docker run -d -p 8081:8081 --name matching-service-container --network matching-service-network matching-service
```

### Stopping Docker Processes and Network

```bash
# To stop the container for matching service, use the container ID for matching service from the previous command
docker stop <container_id_for_matching_service>

# To remove the container, use the container ID for matching service from the previous command
docker remove <container_id_for_matching_service>

# To stop the container for redis, use the container ID for redis from the previous command
docker stop <container_id_for_redis>

# To stop the container for redis, use the container ID for redis from the previous command
docker remove <container_id_for_redis>

# To remove the matching service network, use the following command
docker network rm matching-service-network
```
