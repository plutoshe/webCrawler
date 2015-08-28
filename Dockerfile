FROM golang
MAINTAINER Pluto.She <plutoshe@gmail.com>
ADD . /go/src/main
WORKDIR /go/src/main
RUN go get
RUN go build -o cralwer crawler.go
CMD ["/go/src/main/cralwer"]
