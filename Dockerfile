FROM golang:1.16.10-alpine

WORKDIR /fun-platform
ADD . /fun-platform
RUN go build

RUN apk add --no-cache bash
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

ENTRYPOINT ["./fun-platform"]