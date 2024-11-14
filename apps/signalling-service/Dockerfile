FROM node:18-alpine AS base

# Install dependencies only when needed
FROM base AS deps
# Check https://github.com/nodejs/docker-node/tree/b4117f9333da4138b03a546ec926ef50a31506c3#nodealpine to understand why libc6-compat might be needed.
RUN apk add --no-cache libc6-compat
WORKDIR /app

COPY package*.json pnpm-lock.yaml* ./
RUN \
    if [ -f pnpm-lock.yaml ]; then corepack enable pnpm && pnpm i --frozen-lockfile; \
    else echo "Lockfile not found." && exit 1; \
    fi

COPY . .

# Expose port 3001 so it can be mapped by Docker daemon.
EXPOSE 4444

# Define the command to run your app using CMD which defines your runtime.
CMD [ "node", "server.js" ]