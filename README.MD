# transmissionRPC

### transmission RPC API client 

Initial transmission-daemon config: [RPC config](https://github.com/transmission/transmission/wiki/Editing-Configuration-Files#rpc)

Usage example:
```go
const srvAddr = "http://{SRV_IP}:9091/transmission/rpc"
const user = "{USER}"
const pass = "{PASS}"

var torrent t.Transmission
client, err := t.NewClient(srvAddr, user, pass)
panicOnErr(err)

torrent.Http = client
torrent.Paused = true
torrent.DownloadDir = "/tmp"
torrent.Debug = false

addRes, err := torrent.AddFile("./CentOS-8-x86_64-1905-boot.torrent")
panicOnErr(err)
fmt.Printf("%+v\n", addRes)
// {TorrentAdded:{HashString:db1327d2a23c11aeab5b946b1a498fadf6422b49 ID:9 Name:CentOS-8-x86_64-1905-boot}}

allRes, err := torrent.All()
panicOnErr(err)
fmt.Printf("%+v\n", allRes)
// {Torrents:[{ActivityDate:0 AddedDate:0 BandwidthPriority:0 Comment:CentOS x86_64 1905 ISO Error:0 ErrorString: Eta:-1 ID:9 IsFinished:false LeftUntilDone:559941692 and many more ...

oneRes, err := torrent.ByID(addRes.TorrentAdded.ID)
panicOnErr(err)
fmt.Printf("%+v\n", oneRes)
// {Torrents:[{ActivityDate:0 AddedDate:1578182554 BandwidthPriority:0 Comment: Error:0 ErrorString: Eta:-1 ID:0 IsFinished:false and many more ...
```