package proxy

import (
	"context"
	"google.golang.org/grpc"
)
type Director func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error)

func GetDirector() Director {
	return func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		// 该函数负责从 metadata 获取信息 再进行 服务发现 / 负责均衡
		panic("")
	}
}
