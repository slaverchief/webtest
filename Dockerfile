FROM golang:1.24.2-alpine AS builder



WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/webtest

FROM alpine:latest

COPY --from=builder /bin/webtest /bin/webtest

COPY app.env .

CMD ["/bin/webtest"]