package proxy

import (
	"fmt"
	"google.golang.org/grpc/encoding"
	"github.com/golang/protobuf/proto"
)

// 构造了一个 encoding.Codec 类型的实例，
// 该 codec 尝试将 gRPC 消息当作 raw bytes 来实现，当尝试失败后，会有 protoCodec 作为一个后退的 codec
func Codec() encoding.Codec {
	return &rawCodec{&protoCodec{}}
}

type rawCodec struct {
	parentCodec encoding.Codec
}

type frame struct {
	payload []byte
}

// protoCodec 实现 protobuf 的默认的 codec
type protoCodec struct{}

func (protoCodec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}
func (protoCodec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
func (protoCodec) Name() string {
	return "proto"
}

// 序列化函数，
// 尝试将消息转换为*frame类型，并返回frame的payload实现序列化
// 若失败，则采用变量parentCodec中的Marshal进行序列化
func (c *rawCodec) Marshal(v interface{}) ([]byte, error) {
	out, ok := v.(*frame)
	if !ok {
		return c.parentCodec.Marshal(v)
	}
	return out.payload, nil

}

// 反序列化函数，
// 尝试通过将消息转为*frame类型，提取出payload到[]byte，实现反序列化
// 若失败，则采用变量parentCodec中的Unmarshal进行反序列化
func (c *rawCodec) Unmarshal(data []byte, v interface{}) error {
	dst, ok := v.(*frame)
	if !ok {
		return c.parentCodec.Unmarshal(data, v)
	}
	dst.payload = data
	return nil
}
func (c *rawCodec) Name() string {
	return fmt.Sprintf("proxy>%s", c.parentCodec.Name())
}
