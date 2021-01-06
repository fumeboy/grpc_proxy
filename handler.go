package proxy

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type handler struct {
	director Director
}

var clientStreamDescForProxying = &grpc.StreamDesc{
	ServerStreams: true,
	ClientStreams: true,
}

// 该handler以gRPC server的模式来接受数据流，并将受到的数据转发到指定的connection中
func (h *handler) handler(srv interface{}, StreamAnP grpc.ServerStream) (err error) {
	// 获取请求流的目的 Method 名称
	fullMethodName, ok := grpc.MethodFromServerStream(StreamAnP)
	if !ok {
		return status.Errorf(codes.Internal, "failed to get method from server stream")
	}
	// 根据请求流的 meta 信息，判断出正确的对应的目的方
	// 返回一个到目的方的 connection，方便之后实现数据包的透明转发
	outgoingCtx, backendConn, err := h.director(StreamAnP.Context(), fullMethodName)
	if err != nil {
		return err
	}
	defer backendConn.Close()

	// 新发起一个 Stream `Proxy <-> B`
	CtxBnP, CancelBnP := context.WithCancel(outgoingCtx)
	StreamBnP, err := grpc.NewClientStream(CtxBnP, clientStreamDescForProxying, backendConn, fullMethodName)
	if err != nil {
		return err
	}

	// 发送，A->B
	ErrChanA2B := h.forwardA2B(StreamAnP, StreamBnP)
	// 返回，B->A
	ErrChanB2A := h.forwardB2A(StreamBnP, StreamAnP)

	// 数据流结束处理 & 错误处理
	for i := 0; i < 2; i++ {
		select {
		case err = <-ErrChanA2B:
			if err == io.EOF {
				// 正常结束
				StreamBnP.CloseSend()
				break
			} else {
				// 错误处理 (如链接断开、读错误等)
				CancelBnP()
				return status.Errorf(codes.Internal, "failed proxying s2c: %v", err)
			}
		case err = <-ErrChanB2A:
			// 设置 Trailer
			StreamAnP.SetTrailer(StreamBnP.Trailer())
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
	return status.Errorf(codes.Internal, "proxy should never reach this stage.")
}

func (h *handler) forwardA2B(src grpc.ServerStream, dst grpc.ClientStream) chan error {
	ret := make(chan error, 1)
	go func() {
		// *frame即为我们自定义codec中使用到的数据结构
		f := &frame{}
		for {
			if err := src.RecvMsg(f); err != nil {
				ret <- err
				break
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

func (h *handler) forwardB2A(src grpc.ClientStream, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &frame{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err
				break
			}
			if i == 0 {
				// grpc中客户端到服务器的header只能在第一个客户端消息后才可以读取到，
				// 同时又必须在 flush 第一个msg之前写入到流中
				md, err := src.Header()
				if err != nil {
					ret <- err
					break
				}
				if err := dst.SendHeader(md); err != nil {
					ret <- err
					break
				}
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}
