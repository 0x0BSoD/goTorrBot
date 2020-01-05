package transmissionRPC

import (
	"encoding/json"
	"fmt"
)

// Request ===============================
type Request struct {
	Method    string       `json:"method"`
	Arguments ReqArguments `json:"arguments,omitempty"`
	Tag       int          `json:"tag,omitempty"`
	Format    string       `json:"format,omitempty"`
}

type ReqArguments struct {
	Fields            []string `json:"fields,omitempty"`
	IDs               []int    `json:"ids,omitempty"`
	FileName          string   `json:"filename,omitempty"`
	DownloadDir       string   `json:"download-dir,omitempty"`
	MetaInfo          string   `json:"metainfo,omitempty"`
	Paused            bool     `json:"paused,omitempty"`
	PeerLimit         int      `json:"peer-limit,omitempty"`
	BandwidthPriority int      `json:"bandwidth-priority,omitempty"`
	FilesWanted       []string `json:"files-wanted,omitempty"`
	FilesUnwanted     []string `json:"files-unwanted,omitempty"`
	PriorityHigh      []string `json:"priority-high,omitempty"`
	PriorityLow       []string `json:"priority-low,omitempty"`
	PriorityNormal    []string `json:"priority-normal,omitempty"`
	DeleteLocalData   bool     `json:"delete-local-data"`
	Path              string   `json:"path"`
}

// Response ===============================
type Response struct {
	Result    string                 `json:"result"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Tag       int                    `json:"tag,omitempty"`
}

type Torrent struct {
	Torrents []struct {
		ActivityDate      int            `json:"activityDate,omitempty"`
		AddedDate         int            `json:"addedDate,omitempty"`
		BandwidthPriority int            `json:"bandwidthPriority"`
		Comment           string         `json:"comment,omitempty"`
		Error             int            `json:"error,omitempty"`
		ErrorString       string         `json:"errorString,omitempty"`
		Eta               int            `json:"eta,omitempty"`
		ID                int            `json:"id,omitempty"`
		IsFinished        bool           `json:"isFinished,omitempty"`
		LeftUntilDone     int            `json:"leftUntilDone,omitempty"`
		Name              string         `json:"name,omitempty"`
		PercentDone       float64        `json:"percentDone,omitempty"`
		SizeWhenDone      int            `json:"sizeWhenDone,omitempty"`
		StartDate         int            `json:"startDate,omitempty"`
		Status            int            `json:"status,omitempty"`
		TotalSize         int            `json:"totalSize,omitempty"`
		UploadRatio       float64        `json:"uploadRatio,omitempty"`
		Peers             []ArgPeers     `json:"peers,omitempty"`
		RateDownload      int            `json:"rateDownload,omitempty"`
		RateUpload        int            `json:"rateUpload,omitempty"`
		Files             []ArgFiles     `json:"files,omitempty"`
		FileStats         []ArgFileStats `json:"fileStats,omitempty"`
	} `json:"torrents"`
}

type ArgFiles struct {
	BytesCompleted int    `json:"bytesCompleted"`
	Length         int    `json:"length"`
	Name           string `json:"name"`
}

type ArgFileStats struct {
	BytesCompleted int  `json:"bytesCompleted"`
	Wanted         bool `json:"wanted"`
	Priority       int  `json:"priority"`
}

type ArgPeers struct {
	Address            string  `json:"address"`
	ClientName         string  `json:"clientName"`
	ClientIsChoked     bool    `json:"clientIsChoked"`
	ClientIsInterested bool    `json:"clientIsInterested"`
	FlagStr            string  `json:"flagStr"`
	IsDownloadingFrom  bool    `json:"isDownloadingFrom"`
	IsEncrypted        bool    `json:"isEncrypted"`
	IsIncoming         bool    `json:"isIncoming"`
	IsUploadingTo      bool    `json:"isUploadingTo"`
	IsUTP              bool    `json:"isUTP"`
	PeerIsChoked       bool    `json:"peerIsChoked"`
	PeerIsInterested   bool    `json:"peerIsInterested"`
	Port               int     `json:"port"`
	Progress           float64 `json:"progress"`
	RateToClient       int     `json:"rateToClient"`
	RateToPeer         int     `json:"rateToPeer"`
}

type Added struct {
	TorrentAdded struct {
		HashString string `json:"hashString,omitempty"`
		ID         int    `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
	} `json:"torrent-added"`
}

type Duplicate struct {
	TorrentDuplicate struct {
		HashString string `json:"hashString,omitempty"`
		ID         int    `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
	} `json:"torrent-duplicate"`
}

// Other ===============================
type Info struct {
	AltSpeedDown              int    `json:"alt-speed-down"`
	AltSpeedEnabled           bool   `json:"alt-speed-enabled"`
	AltSpeedTimeBegin         int    `json:"alt-speed-time-begin"`
	AltSpeedTimeDay           int    `json:"alt-speed-time-day"`
	AltSpeedTimeEnabled       bool   `json:"alt-speed-time-enabled"`
	AltSpeedTimeEnd           int    `json:"alt-speed-time-end"`
	AltSpeedUp                int    `json:"alt-speed-up"`
	BlocklistEnabled          bool   `json:"blocklist-enabled"`
	BlocklistSize             int    `json:"blocklist-size"`
	BlocklistURL              string `json:"blocklist-url"`
	CacheSizeMb               int    `json:"cache-size-mb"`
	ConfigDir                 string `json:"config-dir"`
	DhtEnabled                bool   `json:"dht-enabled"`
	DownloadDir               string `json:"download-dir"`
	DownloadDirFreeSpace      int    `json:"download-dir-free-space"`
	DownloadQueueEnabled      bool   `json:"download-queue-enabled"`
	DownloadQueueSize         int    `json:"download-queue-size"`
	Encryption                string `json:"encryption"`
	IdleSeedingLimit          int    `json:"idle-seeding-limit"`
	IdleSeedingLimitEnabled   bool   `json:"idle-seeding-limit-enabled"`
	IncompleteDir             string `json:"incomplete-dir"`
	IncompleteDirEnabled      bool   `json:"incomplete-dir-enabled"`
	LpdEnabled                bool   `json:"lpd-enabled"`
	PeerLimitGlobal           int    `json:"peer-limit-global"`
	PeerLimitPerTorrent       int    `json:"peer-limit-per-torrent"`
	PeerPort                  int    `json:"peer-port"`
	PeerPortRandomOnStart     bool   `json:"peer-port-random-on-start"`
	PexEnabled                bool   `json:"pex-enabled"`
	PortForwardingEnabled     bool   `json:"port-forwarding-enabled"`
	QueueStalledEnabled       bool   `json:"queue-stalled-enabled"`
	QueueStalledMinutes       int    `json:"queue-stalled-minutes"`
	RenamePartialFiles        bool   `json:"rename-partial-files"`
	RPCVersion                int    `json:"rpc-version"`
	RPCVersionMinimum         int    `json:"rpc-version-minimum"`
	ScriptTorrentDoneEnabled  bool   `json:"script-torrent-done-enabled"`
	ScriptTorrentDoneFilename string `json:"script-torrent-done-filename"`
	SeedQueueEnabled          bool   `json:"seed-queue-enabled"`
	SeedQueueSize             int    `json:"seed-queue-size"`
	SeedRatioLimit            int    `json:"seedRatioLimit"`
	SeedRatioLimited          bool   `json:"seedRatioLimited"`
	SpeedLimitDown            int    `json:"speed-limit-down"`
	SpeedLimitDownEnabled     bool   `json:"speed-limit-down-enabled"`
	SpeedLimitUp              int    `json:"speed-limit-up"`
	SpeedLimitUpEnabled       bool   `json:"speed-limit-up-enabled"`
	StartAddedTorrents        bool   `json:"start-added-torrents"`
	TrashOriginalTorrentFiles bool   `json:"trash-original-torrent-files"`
	Units                     struct {
		MemoryBytes int      `json:"memory-bytes"`
		MemoryUnits []string `json:"memory-units"`
		SizeBytes   int      `json:"size-bytes"`
		SizeUnits   []string `json:"size-units"`
		SpeedBytes  int      `json:"speed-bytes"`
		SpeedUnits  []string `json:"speed-units"`
	} `json:"units"`
	UtpEnabled bool   `json:"utp-enabled"`
	Version    string `json:"version"`
}

type Statistics struct {
	ActiveTorrentCount int `json:"activeTorrentCount"`
	CumulativeStats    struct {
		DownloadedBytes int `json:"downloadedBytes"`
		FilesAdded      int `json:"filesAdded"`
		SecondsActive   int `json:"secondsActive"`
		SessionCount    int `json:"sessionCount"`
		UploadedBytes   int `json:"uploadedBytes"`
	} `json:"cumulative-stats"`
	CurrentStats struct {
		DownloadedBytes int `json:"downloadedBytes"`
		FilesAdded      int `json:"filesAdded"`
		SecondsActive   int `json:"secondsActive"`
		SessionCount    int `json:"sessionCount"`
		UploadedBytes   int `json:"uploadedBytes"`
	} `json:"current-stats"`
	DownloadSpeed      int `json:"downloadSpeed"`
	PausedTorrentCount int `json:"pausedTorrentCount"`
	TorrentCount       int `json:"torrentCount"`
	UploadSpeed        int `json:"uploadSpeed"`
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
