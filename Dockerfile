FROM golang:1.18-alpine3.15
RUN mkdir api
COPY . /post-service
WORKDIR /post-service
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 8020
