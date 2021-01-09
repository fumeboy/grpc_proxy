package proxy

var nodes = map[string]string{
	"guestbook": "127.0.0.1:9100",
}// 专为 demo 用的临时的 registry，如果要加上 etcd 的话，作为示例来说就太大了

func fakeBalancer(nodeName string) string {
	return nodes[nodeName]
}
