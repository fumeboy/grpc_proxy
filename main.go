package proxy

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

func proxy(){
	// 无论是请求方的数据流，还是响应方的数据流，对于proxy服务来说都是数据流的进入
	// 也即是proxy需要作为一个server的身份来处理这些请求
	s := grpc.NewServer()

	// 为了实现透明代理，需要自定义特殊的编解码
	encoding.RegisterCodec(Codec())
}