FROM golang:1.21.4 as app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal internal
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/veun-http-demo ./cmd/demo-server

FROM scratch
COPY --from=app /app/veun-http-demo /go/bin/app
ENTRYPOINT ["/go/bin/app"]
