package proxy

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

var conns = map[string]*grpc.ClientConn{}
var globalContext, cancelGlobalContext = context.WithCancel(context.Background())

func RunProxy(port int){
	defer func() {
		cancelGlobalContext()
		for _,conn := range conns{
			_ = conn.Close()
		}
	}()
	// 无论是请求方的数据流，还是响应方的数据流，对于proxy服务来说都是数据流的进入
	// 也即是 proxy 需要作为一个server的身份来处理这些请求
	s := grpc.NewServer(
		grpc.UnknownServiceHandler(handler),
		// proxy 是一个没有方法的 grpc server， 所以任何请求都是 Unknown
		// 故用 UnknownServiceHandler 处理
	)
	l,_ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	_ = s.Serve(l)
}