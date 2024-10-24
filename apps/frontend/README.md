This is the frontend for the question service.

## Tech Stack

- Next.js
- TypeScript
- Ant Design
- SCSS

## Getting Started

First, install the dependencies:

```bash
npm install -g pnpm

pnpm install --frozen-lockfile

# if pnpm install --frozen-lockfile fails, try running
pnpm install
```

Then, follow the `.env.example` file and create a `.env` file in the current directory. Replace the necessary values within.

```bash
NEXT_PUBLIC_QUESTION_SERVICE_URL="http://localhost:8080"
NEXT_PUBLIC_USER_SERVICE_URL="http://localhost:3001/"
NEXT_PUBLIC_MATCHING_SERVICE_URL="ws://localhost:8081"
```

First, run the development server:

```bash
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

## Build Dockerfile

```sh
# Navigate to the frontend app directory
cd apps/frontend

# Build dockerfile (Ensure that your docker daemon is running beforehand)
docker build -t frontend -f Dockerfile .
```

Run the backend server locally and visit http://localhost:3000/ to see the frontend application working

## Running Docker Image

```sh
# Run the docker image, the -d tag is to run it detached
docker run -p 3000:3000 --env-file .env -d frontend

# To see the running container
docker ps

# To stop the container, copy the container id from the previous command
docker stop <container_id>
```
# GCP Deployment Related
## Pushing Docker Image onto Docker Hub
You can push a new image to the Docker Hub repository using the CLI:

```bash
docker tag local-image:tagname docker-hub-repo:tagname
docker push docker-hub-repo:tagname
```

Make sure to replace tagname with your desired image repository tag.

## Naming convention for container image URL in Google Cloud Run
```bash
docker.io/<dockerRepoName>:version
```
Do note that Google Cloud Run caches the previous docker image with the same tag. Therefore, if you want to host the updated docker image, you MUST change the version tag. Otherwise, Google Cloud Run will continuously deploy from the old tag, despite the fact that the docker image has been updated

## How to deploy docker image to Google Cloud Run
Before running the below commands, make sure that you have changed the environmental variables in .env, so that the frontend can interact with all the services that had been deployed on Google Cloud Platform
```
NEXT_PUBLIC_QUESTION_SERVICE_URL="https://cs3219-g24-question-service-manual-534153711557.us-central1.run.app/"
NEXT_PUBLIC_USER_SERVICE_URL="https://cs3219-g24-user-service-manual-534153711557.us-central1.run.app/"
NEXT_PUBLIC_MATCHING_SERVICE_URL="ws://localhost:8081/match"
```
1. Build Image
```bash
docker build -t frontend -f Dockerfile .
```
2. Tag Image
```bash
docker tag local-image:tagname docker-hub-repo:tagname
```
3. Push Image to Docker Hub
```bash
docker push docker-hub-repo:tagname
```
4. Deploy Container with Container Image URL
```bash
docker.io/<dockerRepoName>:version
```

## Below are the following URLs for each services
Frontend
[Link](https://cs3219-g24-frontend-manual-534153711557.us-central1.run.app/)

Question Service
[Link](https://cs3219-g24-question-service-manual-534153711557.us-central1.run.app)

User Service
[Link](https://cs3219-g24-user-service-manual-534153711557.us-central1.run.app)
