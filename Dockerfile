FROM golang:1.9.0

RUN go get github.com/Masterminds/glide

WORKDIR $GOPATH/src/github.com/Masterminds/glide

RUN make build
RUN glide -v