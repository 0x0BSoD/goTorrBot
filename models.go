package main

type Parameters struct {
	Method    string    `json:"method"`
	Arguments Arguments `json:"arguments"`
}

type Arguments struct {
	Fields   []string `json:"fields,omitempty"`
	IDs      []int    `json:"ids,omitempty"`
	FileName string   `json:"filename,omitempty"`
}

type Statistics struct {
	Arguments struct {
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
	} `json:"arguments"`
	Result string `json:"result"`
}
