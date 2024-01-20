FROM golang:1.21 as base

FROM base as dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY api/ ./api/
COPY config/ ./config/
COPY wallets/ ./wallets/

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /application

CMD ["/application"]