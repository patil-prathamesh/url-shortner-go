FROM golang:alpine as builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine

RUN adduser -S -D -H -h /app appuser

USER appuser

WORKDIR /app

COPY --from=builder /build/main .

EXPOSE 3000

CMD [ "./main" ]