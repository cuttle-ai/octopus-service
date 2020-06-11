FROM golang:1.13
LABEL maintainer="Melvin Davis <melvinodsa@gmail.com>"

WORKDIR /app

COPY octopus-service/go.mod octopus-service/go.sum ./
COPY auth-service/go.mod auth-service/go.sum ./
COPY go-sdk/go.mod go-sdk/go.sum ./
COPY db-toolkit/go.mod db-toolkit/go.sum ./
COPY configs/go.mod configs/go.sum ./
COPY octopus/go.mod octopus/go.sum ./
COPY brain/go.mod brain/go.sum ./
COPY auth-service /auth-service
COPY brain /brain
COPY octopus /octopus
COPY configs /configs
COPY db-toolkit /db-toolkit
COPY go-sdk /go-sdk


# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install ca-certificates
RUN apt-get update && apt-get install ca-certificates

# Copy the source from the current directory to the Working Directory inside the container
COPY octopus-service /app

# Command to run the executable
CMD ["go", "run", "main.go"]