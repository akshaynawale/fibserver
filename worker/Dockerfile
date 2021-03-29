FROM golang:1.14

# get the required libs which we need to run  the 
RUN  go get -v "google.golang.org/grpc"

WORKDIR /go/src/fibServer
# copy  all files to this folder
COPY worker worker
# Also copy the protobuff files
COPY proto proto

WORKDIR /go/src/fibServer/worker

RUN go build -o workerexe

CMD ["./workerexe"]