FROM golang:1.17 as build

WORKDIR /build

COPY go.* .

RUN go mod download
COPY . .
#  go mod tidy -compat=1.17
# go build -o app
RUN go build -mod=readonly -v -o app
