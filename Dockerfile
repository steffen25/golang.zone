# golang.zone
FROM golang:alpine
RUN apk update
RUN apk add bash
ADD golang-zone-dump.sql /golang-zone-dump.sql
ADD api api
ADD config/app.json config/app.json
RUN ls -a
EXPOSE 8080
WORKDIR /go
ENTRYPOINT ["./api"]