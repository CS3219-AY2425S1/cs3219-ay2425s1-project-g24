# PeerPrep Docker Compose Guide

This project uses Docker Compose to manage multiple services such as a frontend, backend, and a database. The configuration is defined in the `docker-compose.yml` file, and environment variables can be stored in environment files for different environments (e.g., development, production).

## Prerequisites

Before you begin, ensure you have the following installed on your machine:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Project Structure

In the `./apps` directory:

```plaintext
.
├── docker-compose.yml       # Docker Compose configuration
├── .env                     # Global environment variables (optional)
├── frontend
│   ├── Dockerfile           # Dockerfile for frontend
│   └── ... (other frontend files)
├── question-service
│   ├── Dockerfile           # Dockerfile for question-service
│   └── ... (other question-service files)
├── user-service
│   ├── Dockerfile           # Dockerfile for user-service
│   └── ... (other user-service files)
└── README.md                # Project documentation (for docker compose)
```

## Docker Compose Setup

By using multiple Dockerfiles in Docker Compose, we can manage complex multi-container applications where each service has its own environment and build process.

1. Build and Start the Application

To build and run both the frontend and backend services, you can change your directory to the `./apps` directory and run:

```bash
docker-compose up --build
```

This will:

- Build the Docker images for all services using the specified Dockerfiles
- Start the containers and map the defined ports

2. Access the Application

Once running, you can access:

- The **frontend** at http://localhost:3000
- The **user service** at http://localhost:3001
- The **question service** at http://localhost:8080

3. Stopping Services

To stop the running services, run:

```bash
docker-compose down
```

This command will stop and remove the containers, networks, and volumes created by docker-compose up.

## Troubleshooting

**Common Issues**

- Port Conflicts: If you encounter port conflicts, ensure the host ports specified in docker-compose.yml (e.g., 3000:3000) are not in use by other applications.
- Environment Variables Not Loaded: Ensure the `.env` files are in the correct directories as found in the `docker-compose.yml` file.

**Logs**

You can view the logs for each service using the following command:

```bash
docker-compose logs
```

**Useful Commands**

Rebuild a specific service:

```bash
docker-compose build <service_name>
```

Start services in detached mode (run in the background):

```bash
docker-compose up -d
```

Remove all containers, networks, and volumes created by Docker Compose:

```bash
docker-compose down --volumes
```
