package main

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
}

// Response ===============================
type Response struct {
	Result    string                 `json:"result"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Tag       int                    `json:"tag,omitempty"`
}

type Torrent struct {
	Torrents []struct {
		ActivityDate  int           `json:"activityDate,omitempty"`
		AddedDate     int           `json:"addedDate,omitempty"`
		Comment       string        `json:"comment,omitempty"`
		Error         int           `json:"error,omitempty"`
		ErrorString   string        `json:"errorString,omitempty"`
		Eta           int           `json:"eta,omitempty"`
		ID            int           `json:"id,omitempty"`
		IsFinished    bool          `json:"isFinished,omitempty"`
		LeftUntilDone int           `json:"leftUntilDone,omitempty"`
		Name          string        `json:"name,omitempty"`
		PercentDone   float64       `json:"percentDone,omitempty"`
		SizeWhenDone  int           `json:"sizeWhenDone,omitempty"`
		StartDate     int           `json:"startDate,omitempty"`
		Status        int           `json:"status,omitempty"`
		TotalSize     int           `json:"totalSize,omitempty"`
		UploadRatio   int           `json:"uploadRatio,omitempty"`
		Peers         []interface{} `json:"peers,omitempty"`
		RateDownload  int           `json:"rateDownload,omitempty"`
		RateUpload    int           `json:"rateUpload,omitempty"`
		Files         []ArgFiles    `json:"files,omitempty"`
	} `json:"torrents"`
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

type ArgFiles struct {
	BytesCompleted int    `json:"bytesCompleted"`
	Length         int    `json:"length"`
	Name           string `json:"name"`
}

// Other ===============================
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
