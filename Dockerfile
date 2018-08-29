FROM golang:1.11-alpine
COPY . /go/src/app
WORKDIR /go/src/app
RUN go build
CMD /go/src/app/clustermanager 9292
