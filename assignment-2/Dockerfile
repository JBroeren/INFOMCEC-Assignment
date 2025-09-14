FROM golang:latest

RUN apt-get update && \
    apt install -y curl

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

RUN mkdir -p /usr/src/app/auth

# maunt auth folder as a volume
VOLUME ["/usr/src/app/auth"]

# copy go script
COPY script.go .

ENTRYPOINT ["go", "run", "script.go"]