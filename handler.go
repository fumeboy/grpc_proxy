package proxy

import (
	"github.com/fumeboy/grpc-go"
	"github.com/fumeboy/grpc-go/codes"
	"github.com/fumeboy/grpc-go/status"
	"time"
)

// 该handler以gRPC server的模式来接受数据流，并将受到的数据转发到指定的connection中
func handler(_ interface{}, StreamAnP grpc.ServerStream) (err error) {
	// 获取请求流的目的 Method 名称
	fullMethodName, ok := grpc.MethodFromServerStream(StreamAnP)
	if !ok {
		return status.Errorf(codes.Internal, "failed to get method from server stream")
	}

	// 根据请求流头部信息，判断出正确的对应的目的方
	// 返回一个到目的方的 ip addr
	endpoint, err := director(fullMethodName)
	if err != nil {
		return err
	}

	var conn *grpc.ClientConn
	if conn,ok = conns[endpoint]; !ok { // conn 复用
		conn, err = grpc.DialContext(globalContext, endpoint, grpc.WithInsecure(), grpc.WithTimeout(10*time.Millisecond))
		if err != nil {
			return err
		}
		conns[endpoint] = conn
	}

	err = StreamAnP.(grpc.ServerStreamRedirect).Redirect(conn, fullMethodName)
	if err != nil {
		panic(err)
	}
	return nil
}
