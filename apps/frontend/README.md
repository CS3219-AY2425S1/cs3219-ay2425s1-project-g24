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
