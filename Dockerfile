FROM golang:1.20
RUN mkdir app
RUN chmod 777 -R ./app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fetch ./cmd/fetch/main.go
RUN chmod +x ./fetch

ENTRYPOINT ["./fetch"]