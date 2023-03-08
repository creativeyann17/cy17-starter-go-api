FROM golang:1.19-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o api

FROM alpine

WORKDIR /

COPY --from=build /app/api .

EXPOSE 8080

ENTRYPOINT ["/api"]