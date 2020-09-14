FROM golang:latest

WORKDIR /home/a/go/src/shkabatur/test-service/

COPY . .

RUN go get -d -v ./...

RUN go build .

EXPOSE 8080

CMD ["./test-service"]