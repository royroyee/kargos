FROM golang:1.19

RUN mkdir /agent
WORKDIR /agent

COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /build

ENTRYPOINT go run main.go