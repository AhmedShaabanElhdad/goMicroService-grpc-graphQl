FROM golang:1.23-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates git

WORKDIR /go/src/github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl

COPY go.mod go.sum ./
RUN go mod download

COPY account account
RUN go build -o /go/bin/app ./account/cmd/account

# final image
FROM alpine:3.19
WORKDIR /usr/bin
COPY --from=build /go/bin/app .
EXPOSE 8080
CMD ["app"]
