FROM golang:1.23.2 AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o /bin/todo-api

FROM debian:buster-slim

COPY --from=build /bin/todo-api /bin

EXPOSE 1323

CMD [ "/bin/todo-api" ]