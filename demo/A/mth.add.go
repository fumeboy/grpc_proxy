package main

import (
	"context"
	"github.com/fumeboy/grpc-go"
	"grpc_proxy/demo/A/guestbook"
)

func Add(ctx context.Context, r *guestbook.AddRequest) (resp *guestbook.AddResponse, err error) {
	conn, err := grpc.Dial(proxyAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	req := r
	client := guestbook.NewGuestBookServiceClient(conn)
	resp, err = client.Add(ctx, req)
	return
}
