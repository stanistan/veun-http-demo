FROM golang:1.21.6 as app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ENV GOOS=linux
ENV GOARCH=amd64

COPY . .
RUN go generate ./...
RUN go build -o veun-http-demo ./cmd/demo-server

FROM scratch
COPY --from=app /app/veun-http-demo /veun-http-demo

ENTRYPOINT ["/veun-http-demo"]
