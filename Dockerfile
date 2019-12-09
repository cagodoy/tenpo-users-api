# Base container for compile service
FROM golang AS builder

# Go to builder workdir
WORKDIR /go/src/github.com/cagodoy/tenpo-users-api/

# Copy go modules files
COPY go.mod .
COPY go.sum .

# Install dependencies
RUN go mod download

# Copy all source code
COPY . .

# Compile service
RUN cd cmd/server && GOOS=linux GOARCH=amd64 go build -o ../../bin

#####################################################################
#####################################################################

# Base container for run service
FROM alpine

# Go to workdir
WORKDIR /src/tenpo-users-api

# Install dependencies
RUN apk add --update ca-certificates

# Copy binaries
COPY --from=builder /go/src/github.com/cagodoy/tenpo-users-api/bin/tenpo-users-api /usr/bin/tenpo-users-api
COPY bin/goose /usr/bin/goose
COPY bin/wait-db /usr/bin/wait-db
COPY database/migrations/* /src/tenpo-users-api/migrations/

# Expose service port
EXPOSE 5020

# Run service
CMD ["/bin/sh", "-l", "-c", "wait-db && cd /src/tenpo-users-api/migrations/ && goose postgres ${POSTGRES_DNS} up && tenpo-users-api"]