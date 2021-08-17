FROM golang:1.16-alpine3.14

WORKDIR /usr/src/app

# Install watcher 
RUN go get github.com/canthefason/go-watcher && go install github.com/canthefason/go-watcher/cmd/watcher@latest

# Copy go app, can be onverwriten at runtime
COPY . .

