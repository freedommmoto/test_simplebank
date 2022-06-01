#build stage
FROM golang:1.18-alpine3.15 as builder
WORKDIR /app
COPY . .
RUN go get -u ./...
RUN go build
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
RUN chmod +x migrate

#last stage
FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/app.env .
COPY --from=builder /app/test_simplebank .
COPY --from=builder /app/migrate ./migrate
COPY start.sh .
COPY db/migration /app/db/migration
COPY wait-forv2.2.3.sh /app/wait-forv2.2.3.sh
# RUN ./banksystem

EXPOSE 8082
CMD ["/app/test_simplebank"]
ENTRYPOINT [ "/app/start.sh" ]


