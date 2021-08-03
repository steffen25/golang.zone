FROM golang:1.16.6-buster

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/v1.15.1/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

# database migrations
RUN go get -u github.com/pressly/goose/cmd/goose

COPY ./docker/util/wait-for-it.sh /

WORKDIR /code/golang-zone

CMD ["/wait-for-it.sh", "app_mysql:3306", "--", "air"]