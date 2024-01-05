FROM golang:1.21.5-bullseye AS build

RUN apt-get update && apt-get install -y git

WORKDIR /app

RUN echo wishlist-service

RUN git clone https://github.com/akshay0074700747/wishlist-service-grpc.git .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o bin/wishlist-service

COPY /cmd/.env /app/cmd/bin/

FROM busybox:latest

WORKDIR /wishlist-service

COPY --from=build /app/cmd/bin/wishlist-service .

COPY --from=build /app/cmd/bin/.env .

EXPOSE 50007

CMD ["./wishlist-service"]