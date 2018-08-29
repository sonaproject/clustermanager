FROM golang:1.8
WORKDIR /go/src/app
RUN git clone https://github.com/sonaproject/clustermanager.git
WORKDIR /go/src/app/clustermanager
RUN go build
CMD /go/src/app/clustermanager/clustermanager
