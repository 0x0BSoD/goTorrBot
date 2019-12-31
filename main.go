package main

import "fmt"

const srvAddr = "http://--:9091/transmission/rpc"
const user = "--"
const pass = "--"

func main() {
	c, err := NewClient(srvAddr, user, pass)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", c)

	c.TorrentGet()
}
