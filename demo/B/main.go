package main

import (
	"fmt"
	"github.com/fumeboy/grpc-go"
	"grpc_proxy/demo/B/guestbook"
	"net"
)

type server struct{}
var serverIns guestbook.GuestBookServiceServer = &server{}

func main() {
	srv := grpc.NewServer()
	guestbook.RegisterGuestBookServiceServer(srv, serverIns)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9100))
	if err != nil {
		panic("failed launch server")
	}
	fmt.Println("server running")
	srv.Serve(lis)
}
