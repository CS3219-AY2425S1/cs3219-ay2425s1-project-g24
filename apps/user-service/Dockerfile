FROM node:20-alpine AS base

FROM base AS deps

WORKDIR /app

COPY package*.json ./

RUN npm install

# Uncomment the following line if you are building code for production.
# RUN npm ci --omit=dev

COPY . .

# Expose port 3001 so it can be mapped by Docker daemon.
EXPOSE 3001

# Define the command to run your app using CMD which defines your runtime.
CMD [ "npm", "start" ]