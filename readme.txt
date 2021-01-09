# grpc proxy

简单的实现 grpc 代理

主要参考： https://github.com/mwitkow/grpc-proxy
它也是目前找到的唯一的代理实现

## codec

对于 contentSubtype == "" 的消息 (grpc/rpc_util.go/setCallInfoCodec)，我们需要采用默认的 protobuf codec
对于要转发的数据，我们需要一个对协议的无感知仅作转发的 codec
故实现 ./codec.go/rawCodec

## director

服务发现 / 负载均衡

## handler

handler 实现 `A -> P, P -> B` 和 `B -> P, P -> A` 中间的数据转发

# demo

执行
    make init
    make run1
    make run2
    make run3

其中 init 需要你安装了 protoc 和 protoc-gen-go

demo 中的 B 和 P 分别使用端口 9100 和 21000

# TODO

如果看了 handler 的代码， 其实可以发现其中还是有很大的问题
比如 读到的数据还需要多余的一次序列化、反序列化，能不能跳过序列化直接转发？
比如 为了享受 HTTP2 的优势， 还是需要手动管理多个 conns


