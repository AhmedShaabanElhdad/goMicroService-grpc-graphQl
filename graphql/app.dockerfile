FROM golang:1.23-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates git

WORKDIR /go/src/github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN GO111MODULE=on go build -o /go/bin/app ./graphql

FROM alpine:3.19
WORKDIR /usr/bin
COPY --from=build /go/bin/app .
EXPOSE 8081
CMD ["app"]