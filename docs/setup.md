# Setup Guide

## Instructions to set up

### Setting up the local environment

1. Visit the repository at the following [url](https://github.com/CS3219-AY2425S1/cs3219-ay2425s1-project-g24)
2. Clone the repository
3. Ensure that you are in the correct branch (main)
   a. Open terminal and run the following command: `git branch`
4. Install Docker
   a. https://www.docker.com/get-started
5. Install Docker Compose
   a. https://docs.docker.com/compose/install/

### Setting up the database secrets

1. Copy the `cs3219-staging-codeexecution-firebase-adminsdk-ce48j-00ab09514c.json` file found in the zip file into the `./apps/execution-service/` directory.
2. Copy the `cs3219-staging-codehisto-bb61c-firebase-adminsdk-egopb-95cfaf9b87.json` file found in the zip file into the `./apps/history-service/` directory.
3. Copy the `cs3219-g24-firebase-adminsdk-9cm7h-b1675603ab.json` file found in the zip file into the `./apps/question-service/` directory.

### Setting up environment variables

1. Create/Update the `./apps/execution-service/.env` file with the following variables:

```
FIREBASE_CREDENTIAL_PATH=cs3219-staging-codeexecution-firebase-adminsdk-ce48j-00ab09514c.json
PORT=8083
HISTORY_SERVICE_URL=http://history-service:8082/
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
```

2. Create/Update the `./apps/frontend/.env` file with the following variables:

```
NEXT_PUBLIC_QUESTION_SERVICE_URL="http://localhost:8080/"
NEXT_PUBLIC_USER_SERVICE_URL="http://localhost:3001/"
NEXT_PUBLIC_MATCHING_SERVICE_URL="ws://localhost:8081/match"
NEXT_PUBLIC_SIGNALLING_SERVICE_URL="ws://localhost:4444/"
NEXT_PUBLIC_HISTORY_SERVICE_URL="http://localhost:8082/"
NEXT_PUBLIC_EXECUTION_SERVICE_URL="http://localhost:8083/"
```

3. Create/Update the `./apps/history-service/.env` file with the following variables:

```
FIREBASE_CREDENTIAL_PATH=cs3219-staging-codehisto-bb61c-firebase-adminsdk-egopb-95cfaf9b87.json
PORT=8082
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
```

4. Create/Update the `./apps/matching-service/.env` file with the following variables:

```
PORT=8081
MATCH_TIMEOUT=30
JWT_SECRET=63059c735adba274cde40f2b1c0b955842d531b115bb4df3058d769b173dcc78
REDIS_URL=redis-container:6379
QUESTION_SERVICE_GRPC_URL=question-service:50051
```

5. Create/Update the `./apps/question-service/.env` file with the following variables:

```
FIREBASE_CREDENTIAL_PATH=cs3219-g24-firebase-adminsdk-9cm7h-b1675603ab.json
JWT_SECRET=63059c735adba274cde40f2b1c0b955842d531b115bb4df3058d769b173dcc78
EXECUTION_SERVICE_URL="http://execution-service:8083/"
```

6. Create/Update the `./apps/signalling-service/.env` file with the following variables:

```
PORT=4444
JWT_SECRET=63059c735adba274cde40f2b1c0b955842d531b115bb4df3058d769b173dcc78
```

7. Create/Update the `./apps/user-service/.env` file with the following variables:

```
DB_CLOUD_URI=mongodb+srv://admin:admin_user_service@cs3219-user-service.bbmji.mongodb.net/?retryWrites=true&w=majority&appName=cs3219-user-service
DB_LOCAL_URI=mongodb://127.0.0.1:27017/peerprepUserServiceDB
PORT=3001
ENV=PROD
JWT_SECRET=63059c735adba274cde40f2b1c0b955842d531b115bb4df3058d769b173dcc78
```

### Set-up and run ALL Services via Docker Compose

1. Change directory to `./apps/` directory, run: `cd ./apps`
2. To set up and run all services via docker compose, run: `docker compose up --build`

### Set-up and run Docker Container for Execution Service

1. Change directory to `./apps/execution-service`, run: `cd ./apps/execution-service`
2. To set up the docker container for the execution service, run: `docker build -t execution-service .`
3. To run the docker container for the execution service, run: `docker run -p 8083:8083 --env-file .env -d execution-service`

### Set-up and run Docker Container for Frontend

1. Change directory to `./apps/frontend` directory, run: `cd ./apps/frontend`
2. To set up the docker container for the frontend, run: `docker build -t frontend -f Dockerfile .`
3. To run the docker container for the frontend, run: `docker run -p 3000:3000 --env-file .env -d frontend`

### Set-up and run Docker Container for History Service

1. Change directory to `./apps/history-service` directory, run: `cd ./apps/history-service`
2. To set up the docker container for the history service, run: `docker build -t history-service .`
3. To run the docker container for the history service, run: `docker run -p 8082:8082 -d history-service`

### Set-up and run Docker Container for Matching-Service

1. Change directory to `./apps/matching-service` directory, run: `cd ./apps/matching-service`
2. Ensure that `REDIS_URL=redis-container:6379` in .env file for matching-service
3. To set up the go docker container for the matching service, run: `docker build -f Dockerfile -t match-go-app .`
4. To create the docker network for redis and go, run: `docker network create redis-go-network`
5. To start a new Redis container in detached mode using the redis image from Docker Hub, run: `docker run -d --name redis-container --network redis-go-network redis`
6. To run the go docker container for the matching-service, run: `docker run -d -p 8081:8081 --name go-app-container --network redis-go-network match-go-app`

### Set-up and run Docker Container for Question-Service

1. Change directory to `./apps/question-service` directory, run: `cd ./apps/question-service`
2. To set up the docker container for the question service, run: `docker build -t question-service .`
3. To run the docker container for the question service, run: `docker run -p 8080:8080 --env-file .env -d question-service`

### Set-up and run Docker Container for Signalling-Service

1. Change directory to `./apps/signalling-service` directory, run: `cd ./apps/signalling-service`
2. To set up the docker container for the signalling service, run: `docker build -t signalling-service -f Dockerfile .`
3. To run the docker container for the signalling service, run: `docker run -p 4444:4444 --env-file .env -d signalling-service`

### Set-up and run Docker Container for User-Service

1. Change directory to `./apps/user-service` directory, run: `cd ./apps/user-service`
2. To set up the docker container for the user service, run: `docker build -t user-service -f Dockerfile .`
3. To run the docker container for the user service, run: `docker run -p 3001:3001 --env-file .env -d user-service`

## Running the application

1. To check if all the Docker containers are up, run `docker ps`.
2. Navigate to `http://localhost:3000/register` to register to the application.
3. Navigate to `http://localhost:3000/login` to login to the application.
4. Access the application via an admin account with the following login details:
   a. email: `admin@gmail.com`
   b. password: `admin`
5. Access the application via non-admin account with the following login details (via
   incognito mode to test matchmaking/collaboration function):
   a. email: `notadmin@gmail.com`
   b. password: `notadmin`

NOTE: The MongoDB connection might get blocked on NUS Wi-Fi, so might have to use another Wi-Fi.
