FROM golang:1.23

WORKDIR /usr/src/app

# Install Docker CLI
RUN apt-get update && apt-get install -y \
    docker.io \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

RUN go mod tidy && go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/app ./main.go

EXPOSE 8083

CMD ["app"]
