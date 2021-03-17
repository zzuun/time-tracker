FROM golang:1.14.15-alpine3.12 AS build-env

WORKDIR /go/src/time-tracker/

COPY . .

RUN go build -o main main.go

FROM alpine:3.12

WORKDIR /app

COPY --from=build-env /go/src/time-tracker /app/

EXPOSE 8000

CMD ./main /app
