FROM golang

COPY . /go/src/github.com/kravcs/gogo
WORKDIR /go/src/github.com/kravcs/gogo

RUN go get ./
RUN go build

EXPOSE 8000

CMD ["gogo"]