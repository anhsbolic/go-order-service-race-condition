FROM golang:1.16-alpine

RUN apk add git && apk add --no-cache go
RUN go version

WORKDIR /app
COPY . .
ADD ./env.conf /app/.env
ADD ./env.conf /app/test/.env

RUN go build -o order-service-api .

EXPOSE 3000

CMD ["./order-service-api"]