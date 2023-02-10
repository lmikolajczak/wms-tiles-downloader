FROM golang:1.19-alpine

WORKDIR /code

# Pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them
# in subsequent builds if they change.
COPY go.mod go.sum* ./
RUN go mod download && go mod verify

CMD tail -f /dev/null
