FROM golang

COPY . /go/src/github.com/kravcs/weather_v2
WORKDIR /go/src/github.com/kravcs/weather_v2

RUN go get ./
RUN go build

EXPOSE 8000

CMD ["weather_v2"]