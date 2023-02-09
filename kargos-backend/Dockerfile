FROM golang:1.19

# RUN mkdir -p /host_kube

# TODO builder pattern
RUN mkdir /kargos
WORKDIR /kargos

COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /build

EXPOSE 9000
RUN pwd
ENTRYPOINT go run main.go
