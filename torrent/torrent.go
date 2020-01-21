package torrent

import (
	"fmt"
	"github.com/0x0bsod/torrBot/session"
)

type Torrents struct {
	Session  *session.Session `json:"-"`
	Torrents []*Torrent       `json:"torrents"`
}

type Torrent struct {
	ActivityDate      int            `json:"activityDate,omitempty"`
	AddedDate         int            `json:"addedDate,omitempty"`
	BandwidthPriority int            `json:"bandwidthPriority,omitempty"`
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
	StatusString      string         `json:"status_string,omitempty"`
	TotalSize         int            `json:"totalSize,omitempty"`
	UploadRatio       float64        `json:"uploadRatio,omitempty"`
	Peers             []ArgPeers     `json:"peers,omitempty"`
	RateDownload      int            `json:"rateDownload,omitempty"`
	RateUpload        int            `json:"rateUpload,omitempty"`
	Files             []ArgFiles     `json:"files,omitempty"`
	FileStats         []ArgFileStats `json:"fileStats,omitempty"`
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

//

func (t *Torrents) Verify(ID int) error {
	res, err := t.Session.ExtractArgs(&Request{
		Method: "torrent-verify",
		Arguments: ReqArguments{
			IDs: []int{ID},
		},
	})
	if err != nil {
		return err
	}

	if res.Result == "success" {
		return nil
	}

	return fmt.Errorf("request failed")
}

func (t *Torrent) Start(ID int) error {
	res, err := t.makeCall(&Request{
		Method: "torrent-start",
		Arguments: ReqArguments{
			IDs: []int{ID},
		},
	})
	if err != nil {
		return err
	}

	if res.Result == "success" {
		return nil
	}

	return fmt.Errorf("request failed")
}

func (t *Torrent) Stop(ID int) error {
	res, err := t.makeCall(&Request{
		Method: "torrent-stop",
		Arguments: ReqArguments{
			IDs: []int{ID},
		},
	})
	if err != nil {
		return err
	}

	if res.Result == "success" {
		return nil
	}

	return fmt.Errorf("request failed")
}

func (t *Torrent) Remove(ID int, rmLocalData bool) error {
	res, err := t.makeCall(&Request{
		Method: "torrent-remove",
		Arguments: ReqArguments{
			IDs:             []int{ID},
			DeleteLocalData: rmLocalData,
		},
	})
	if err != nil {
		return err
	}

	if res.Result == "success" {
		return nil
	}

	return fmt.Errorf("request failed")
}

func (t *Torrents) ResolveStatus() {
	for _, i := range t.Torrents {
		i.StatusString = TorrentStatus(i.Status)
	}
}
