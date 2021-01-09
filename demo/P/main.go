package main

import proxy "grpc_proxy"

// 127.0.0.1: 21000
func main() {
	proxy.RunProxy(21000)
}
