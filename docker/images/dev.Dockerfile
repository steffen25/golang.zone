FROM golang:1.16.6-buster

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/v1.27.3/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

WORKDIR /code/golang-zone

COPY ./docker/util/wait-for-it.sh /

CMD ["/wait-for-it.sh", "golang_zone_db:5432", "--", "air"]
