FROM golang:1.15.6

RUN go get github.com/Masterminds/glide
RUN go get -u github.com/pressly/goose/cmd/goose

WORKDIR $GOPATH/src/github.com/Masterminds/glide

RUN make build
RUN glide -v