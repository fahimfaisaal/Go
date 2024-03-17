FROM golang:1.22-alpine

WORKDIR /app

COPY ./go.* ./

RUN go mod download

COPY . .

RUN apk update && apk upgrade && apk add --no-cache zsh

ARG RUN=""

ENV TYPE=${RUN}

CMD ["go", "run", "main.go"]
