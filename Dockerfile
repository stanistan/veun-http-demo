FROM golang:1.21.4 as app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o veun-http-demo ./cmd/demo-server

FROM scratch
COPY --from=app /app/veun-http-demo /veun-http-demo

ENTRYPOINT ["/veun-http-demo"]
