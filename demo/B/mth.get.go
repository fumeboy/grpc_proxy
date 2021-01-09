package main

import (
	"context"
	"fmt"
	"github.com/fumeboy/grpc-go/codes"
	"github.com/fumeboy/grpc-go/status"
	"grpc_proxy/demo/B/guestbook"
	"grpc_proxy/demo/B/model"
)


func (this *server) Get(ctx context.Context, r *guestbook.GetRequest) (resp *guestbook.GetResponse, err error) {
	fmt.Println("method: Get")
	offset, limit := r.GetOffset(), r.GetLimit()
	if offset < 0 || limit <= 0 {
		err = status.Errorf(codes.InvalidArgument, "add msg failed")
		return nil, err
	}
	if messages, err := model.GetMsg(offset, limit);err != nil {
		return nil, err
	}else{
		resp = &guestbook.GetResponse{}
		for _, one := range messages {
			msg := &guestbook.Msg{
				Email:   one.Email,
				Content: one.Content,
			}
			resp.Msgs = append(resp.Msgs, msg)
		}
		return resp,nil
	}
}
