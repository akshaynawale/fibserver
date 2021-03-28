FROM golang:1.14

WORKDIR /go/src/fibServer
# copy  all files to this folder
COPY webServer webServer
# Also copy the protobuff files
COPY proto proto

# get the required libs which we need to run  the 

RUN  go get -v "github.com/gomodule/redigo/redis" "google.golang.org/grpc"

WORKDIR /go/src/fibServer/webServer

RUN go build -o executable 

CMD ["./executable"]