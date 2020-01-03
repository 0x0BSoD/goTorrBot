package main

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
	torrent.Paused = true
	torrent.DownloadDir = "/tmp"
	torrent.Debug = false

	a, _ := torrent.SessionInfo()
	PrettyPrint(a)
}
