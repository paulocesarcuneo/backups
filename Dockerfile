FROM golang:latest
RUN apt update -y
RUN apt install -y netcat
WORKDIR /go
COPY . .
RUN mkdir -p /go/lib
ENV GOPATH /go/lib
RUN go build agent.go
RUN go build coordinator.go
RUN mkdir -p /data/
RUN wget http://skateipsum.com/get/3/1/text -O /data/lorem
RUN wget http://skateipsum.com/get/3/1/text -O /data/ipsum
