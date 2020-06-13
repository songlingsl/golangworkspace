package main

import (
"github.com/smallnest/rpcx/server"
"songlingkey/rpxc"

)


func main() {
	s := server.NewServer()
	s.RegisterName("Arith", new(rpxc.Arith), "")
	s.Serve("tcp", ":8972")
}
