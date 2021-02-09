FROM golang:1.15-alpine

RUN apk add --update --no-cache alpine-sdk

WORKDIR /pgxscan

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . ./

CMD go test -v -race -bench=. -benchmem --tags=integration ./...