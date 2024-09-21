FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go mod download && go mod vendor

RUN go build -o main cmd/profile/main.go

CMD ["./main"]

FROM alpine

WORKDIR /build

COPY --from=builder /build/main /build/main
COPY ../../configs/config.yaml config/config.yaml

ENV CONFIG_PATH=config/config.yaml

CMD ["./main"]