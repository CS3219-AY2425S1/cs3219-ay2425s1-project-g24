FROM golang:1.23

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

RUN go mod tidy && go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/app ./main.go

EXPOSE 8080 50051

CMD ["app"]
