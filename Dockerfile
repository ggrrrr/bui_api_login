FROM golang:1.17-alpine as stage
WORKDIR /build
COPY go.* /build
RUN go mod download
#  go mod tidy -compat=1.17

RUN go mod tidy 
COPY . /build

# RUN make
RUN go build -mod=readonly -v -o /build/app main.go

FROM alpine:3.14
WORKDIR /app
RUN  mkdir /app/cql

COPY  cql ./cql

COPY --from=stage /build/app /app
CMD /app/app

