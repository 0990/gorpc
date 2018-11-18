package main

import (
	_ "github.com/davyxu/cellnet/peer/tcp" // 注册TCP Peer
	_ "github.com/davyxu/cellnet/proc/tcp" // 注册TCP Processor
)

const peerAddress = "127.0.0.1:17701"

func main() {

}
