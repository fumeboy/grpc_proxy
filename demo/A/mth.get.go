package main

import (
	"context"
	"github.com/fumeboy/grpc-go"
	"github.com/fumeboy/grpc-go/metadata"
	"grpc_proxy/demo/A/guestbook"
)

func Get(ctx context.Context, r*guestbook.GetRequest)(resp*guestbook.GetResponse, err error){
	conn, err := grpc.Dial(proxyAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	req := r
	client := guestbook.NewGuestBookServiceClient(conn)
	ctx = metadata.NewOutgoingContext(ctx, map[string][]string{
	})
	resp, err = client.Get(ctx, req)
	return
}