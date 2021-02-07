FROM golang:1.15-alpine

WORKDIR /pgxscan

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . ./

RUN apk add --update --no-cache alpine-sdk

CMD go test -v -race --tags=integration ./...