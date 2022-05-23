#build stage
FROM golang:1.18-alpine3.15 as builder
WORKDIR /app
COPY . .
RUN go build

#last stage
FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/banksystem .
# RUN ./banksystem

EXPOSE 8082
CMD ["/app/banksystem"]


