# # Base container
# FROM golang

# # Go to workdir
# WORKDIR /go/src/github.com/cagodoy/tenpo-users-api/

# # Copy go modules files
# COPY go.mod .
# COPY go.sum .

# # Install dependencies
# RUN go mod download

# # Compile service
# RUN cd cmd/server && GOOS=linux GOARCH=amd64 go build -o $(BIN)

FROM alpine

RUN apk add --update ca-certificates

WORKDIR /src/tenpo-users-api

COPY bin/tenpo-users-api /usr/bin/tenpo-users-api
COPY bin/goose /usr/bin/goose
COPY bin/wait-db /usr/bin/wait-db
COPY database/migrations/* /src/tenpo-users-api/migrations/

EXPOSE 5020

CMD ["/bin/sh", "-l", "-c", "wait-db && cd /src/tenpo-users-api/migrations/ && goose postgres ${POSTGRES_DNS} up && tenpo-users-api"]