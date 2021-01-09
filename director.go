package proxy

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// 临时写法
func splitFullName(fullname string) (nodeName, serviceName, mthName string, ok bool) {
	fmt.Println(fullname)
	fullname = fullname[1:]
	strs := strings.Split(fullname, "/")
	if len(strs) != 2 {
		return
	}
	strs2 := strings.Split(strs[0], ".")
	if len(strs2) != 2 {
		return
	}
	return strs2[0], strs2[1], strs[1], true
}


// 该函数负责从 metadata 获取信息 再进行 服务发现 / 负载均衡
func director(fullName string) (string, error) {
	nodeName, _, _, ok := splitFullName(fullName)
	if !ok {
		return "", status.Errorf(codes.Internal, "预料之外的格式")
	}
	endpoint := fakeBalancer(nodeName)
	return endpoint, nil
}
