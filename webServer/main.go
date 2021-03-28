package main

import (
	"context"
	pb "fibServer/proto"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
)

var serverAddr string

type WebServer struct {
	workerClient pb.FibWorkerClient
	redisClient  redis.Conn
}

func NewWebServer() WebServer {
	// create a grpc client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...) // grpc server connection
	if err != nil {
		fmt.Printf("failed to connect to the grpc server with err: %v", err)
		os.Exit(1)
	}
	// create redis client
	// Initialize the redis connection to a redis instance running on your local machine
	rc, err := redis.DialURL("redis://fibserver_redis_1") // instead of localhost use the container name instead
	if err != nil {
		panic(err)
	}

	// create a client for our worker
	client := pb.NewFibWorkerClient(conn)

	// register all the handlers
	http.Handle("/", NewFibHandler(client, rc))

	return WebServer{workerClient: client, redisClient: rc}
}

type FibHandler struct {
	Message string
	client  pb.FibWorkerClient // worker client to request the fibonacci data
	tmpl    *template.Template // html template for sending the replay
	rc      redis.Conn
}

func NewFibHandler(client pb.FibWorkerClient, rc redis.Conn) *FibHandler {
	homeTmpl, err := template.ParseFiles("html/home.html")
	if err != nil {
		fmt.Printf("failed to parse the template: %v\n", err)
		os.Exit(1)
	}
	return &FibHandler{tmpl: homeTmpl, client: client, rc: rc}

}

func (fh *FibHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	fmt.Println("now inside the ServerHTTP")
	// read data from the request
	indexStr := req.FormValue("index")
	// check if the output is availabe with redis
	ans, err := fh.rc.Do("GET", indexStr)
	if err != nil {
		fmt.Printf("failed to get output form redis: %v", err)
	} else if ans != nil {
		fh.Message = fmt.Sprintf("found in redis: fibonacci at index: %s is %s", indexStr, ans)
		// write data to the response
		if err = fh.tmpl.Execute(w, fh.Message); err != nil {
			fmt.Printf("failed to execute the template: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// convert the string to the int
	input, err := strconv.Atoi(indexStr)
	var reply *pb.FibReply
	if err != nil {
		fmt.Printf("failed to convert the input to the int: %v\n", err)
		fh.Message = fmt.Sprintf("invalid input: must be a integer: input: %s", indexStr)
	} else {

		// make call to the worker
		request := pb.FibRequest{Num: int32(input)}
		reply, err = fh.client.GetFibNumber(context.Background(), &request)
		if err != nil {
			fmt.Printf("failed to get fib seq err : %v", err)
		}

		fh.Message = fmt.Sprintf("fibonacci number at index %d is %d", input, reply.Num)
	}

	// write data to the response
	if err = fh.tmpl.Execute(w, fh.Message); err != nil {
		fmt.Printf("failed to execute the template: %v\n", err)
		os.Exit(1)
	}
	// store the output to the redis store
	//fmt.Println("setting value")
	_, err = fh.rc.Do("SET", indexStr, fmt.Sprintf("%d", reply.Num))
	if err != nil {
		fmt.Println("failed to set the key: %s value: %d in the redis cache", indexStr, reply.Num)
	}
}

func main() {
	flag.StringVar(&serverAddr, "serverAddr", "fibserver_worker_1:5544", "Fib Worker server address to connect")
	flag.Parse()
	NewWebServer()

	fmt.Println("starting the webserver now")
	//go func() {
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
	//}()
}
