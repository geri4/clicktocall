FROM golang:1.8
MAINTAINER Andrey Gerasimov <grsanw@gmail.com>
WORKDIR /go/src/clicktocall
COPY . .
RUN go-wrapper download
RUN go-wrapper install
EXPOSE 9090
CMD ["go-wrapper", "run", "clicktocall"]
