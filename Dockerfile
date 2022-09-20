FROM golang:1.18-alpine
WORKDIR /app/
COPY . .
RUN apk --update --no-cache add postgresql-client
RUN apk --update --no-cache add make
EXPOSE 8080

CMD ["./migrate.sh", "db"]