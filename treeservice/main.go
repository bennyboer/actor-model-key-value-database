package main

import (
	"flag"
	"fmt"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	printHeader()

	bind := flag.String(
		"bind",
		fmt.Sprintf("%s:%d", defaultName, defaultPort),
		fmt.Sprintf("to what name and address to bind the service: for example --bind=\"%s:%d\"", defaultName, defaultPort),
	)
	flag.Parse()

	listener, e := net.Listen("tcp", *bind)
	if e != nil {
		log.Fatalf("failed to listen: %v", e)
	}

	server := grpc.NewServer()
	messages.RegisterTreeServiceServer(server, &Service{})
	if e := server.Serve(listener); e != nil {
		log.Fatalf("failed to serve: %v", e)
	}
}

func printHeader() {
	fmt.Printf("%s\n\n", welcomeHeader)
}
