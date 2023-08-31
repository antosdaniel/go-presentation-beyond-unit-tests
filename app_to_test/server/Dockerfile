FROM golang:1.21

WORKDIR /src/
COPY . .

RUN go mod vendor

CMD ["go", "run", "app_to_test/server/main.go"]
