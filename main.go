package main

import "fmt"

const srvAddr = "http://--:9091/transmission/rpc"
const user = "--"
const pass = "--"

func main() {
	client, err := NewClient(srvAddr, user, pass)
	if err != nil {
		panic(err)
	}

	var torrent Transmission
	torrent.http = client

	a, _ := torrent.SessionStats()
	PrettyPrint(a)
	fmt.Println("=======")

	//b, _ := torrent.AddMagnet()
	//PrettyPrint(b)
	//fmt.Println("=======")
	//

	c, err := torrent.All()
	if err != nil {
		panic(err)
	}
	PrettyPrint(c)
	fmt.Println("=======")

	//d, _ := torrent.Start(1)
	//PrettyPrint(d)
	//fmt.Println("=======")

	e, _ := torrent.ByID(1)
	PrettyPrint(e)
	fmt.Println("=======")
}
