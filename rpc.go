package transmissionRPC

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// https://github.com/transmission/transmission/blob/master/extras/rpc-spec.txt

type Transmission struct {
	http        *client
	DownloadDir string
	Paused      bool
	Debug       bool
}

// ===========================================
// NewClient return client instance with token
func NewClient(url, user, password string) (*Transmission, error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := client{
		Url:      url,
		User:     user,
		Password: password,
		Http:     &http.Client{Transport: tr},
	}

	resp, err := client.getResponse("GET", "/", nil)
	if err != nil {
		return nil, fmt.Errorf("error during getting token: %s", err)
	}

	client.Token = resp.Header.Get("X-Transmission-Session-Id")

	return &Transmission{http: &client}, nil
}

// =====================================================================================================================
// Get
// =====================================================================================================================

func (t *Transmission) ByID(ID int) ([]*Torrent, error) {
	res, err := t.makeCall(&Request{
		Method: "torrent-get",
		Arguments: ReqArguments{
			Fields: FieldList(Name, Status, ErrorString,
				AddedDate, Peers, IsFinished, LeftUntilDone,
				PercentDone, Eta, TotalSize, RateDownload,
				RateUpload, UploadRatio, Files, FileStats),
			IDs: []int{ID},
		},
	})
	if err != nil {
		return []*Torrent{}, err
	}

	if res.Result == "success" {
		var r Torrents
		err := t.extractArgs(res, &r)
		if err != nil {
			return []*Torrent{}, err
		}
		if len(r.Torrents) == 0 {
			return r.Torrents, fmt.Errorf("no torrents")
		}
		t.resolveStatus(r.Torrents)
		return r.Torrents, nil
	}

	return []*Torrent{}, fmt.Errorf("request failed")
}

func (t *Transmission) ByIDFields(ID int, f ...GetField) ([]*Torrent, error) {
	res, err := t.makeCall(&Request{
		Method: "torrent-get",
		Arguments: ReqArguments{
			Fields: FieldList(f...),
			IDs:    []int{ID},
		},
	})
	if err != nil {
		return []*Torrent{}, err
	}

	if res.Result == "success" {
		var r Torrents
		err := t.extractArgs(res, &r)
		if err != nil {
			return []*Torrent{}, err
		}
		if len(r.Torrents) == 0 {
			return r.Torrents, fmt.Errorf("no torrents")
		}

		t.resolveStatus(r.Torrents)

		return r.Torrents, nil
	}

	return []*Torrent{}, fmt.Errorf("request failed")
}

// =====================================================================================================================
// Add
// =====================================================================================================================

func (t *Transmission) AddFile(path string) (Added, error) {
	// open file and encode to base64
	f, err := os.Open(path)
	if err != nil {
		return Added{}, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	base64Str := base64.StdEncoding.EncodeToString(content)
	// ===

	p := &Request{
		Method: "torrent-add",
		Arguments: ReqArguments{
			MetaInfo: base64Str,
			Paused:   t.Paused,
		},
	}
	if t.DownloadDir != "" {
		p.Arguments.DownloadDir = t.DownloadDir
	}

	res, err := t.makeCall(p)
	if err != nil {
		return Added{}, err
	}

	if res.Result == "success" {
		var r Added
		err := t.extractArgs(res, &r)
		if err != nil {
			return Added{}, err
		}

		if (Added{}) == r {
			var d Duplicate
			err := t.extractArgs(res, &d)
			if err != nil {
				return Added{}, err
			}
			errB, _ := json.Marshal(d)
			return Added{}, fmt.Errorf(string(errB))
		}

		return r, nil
	}

	return Added{}, fmt.Errorf("request failed")
}

func (t *Transmission) AddMagnet(magnetLink string) (Added, error) {
	p := &Request{
		Method: "torrent-add",
		Arguments: ReqArguments{
			FileName: magnetLink,
			Paused:   t.Paused,
		},
	}
	if t.DownloadDir != "" {
		p.Arguments.DownloadDir = t.DownloadDir
	}

	res, err := t.makeCall(p)
	if err != nil {
		return Added{}, err
	}

	if res.Result == "success" {
		var r Added
		err := t.extractArgs(res, &r)
		if err != nil {
			return Added{}, err
		}

		if (Added{}) == r {
			var d Duplicate
			err := t.extractArgs(res, &d)
			if err != nil {
				return Added{}, err
			}
			errB, _ := json.Marshal(d)
			return Added{}, fmt.Errorf(string(errB))
		}

		return r, nil
	}

	return Added{}, fmt.Errorf("request failed")
}

// =====================================================================================================================
// Set
// =====================================================================================================================

const (
	Low int = iota
	Normal
	High
)

// SetPriority gets torrent ID, fileID and level
// if fileID == -1 then all files is be affected
// level const: Low, Normal, High
func (t *Transmission) SetPriority(ID, fileID, level int) error {

	var fileIDs []int

	if fileID == -1 {
		d, err := t.ByIDFields(ID, Files)
		if err != nil {
			return err
		}
		filesCount := len(d[0].Files)
		fileIDs = make([]int, 0, filesCount)

		for i := 0; i < filesCount; i++ {
			fileIDs = append(fileIDs, i)
		}
	} else {
		fileIDs = []int{fileID}
	}

	p := &Request{
		Method:    "torrent-set",
		Arguments: ReqArguments{IDs: []int{ID}},
	}

	switch level {
	case 0:
		p.Arguments.PriorityLow = fileIDs
		break
	case 1:
		p.Arguments.PriorityNormal = fileIDs
		break
	case 2:
		p.Arguments.PriorityHigh = fileIDs
		break
	default:
		return fmt.Errorf("request failed")
	}

	res, err := t.makeCall(p)
	if err != nil {
		return err
	}

	if res.Result == "success" {
		return nil
	}

	return fmt.Errorf("request failed")
}

// =====================================================================================================================
// Other
// =====================================================================================================================

func (t *Transmission) SessionStats() (Statistics, error) {
	res, err := t.makeCall(&Request{
		Method: "session-stats",
	})
	if err != nil {
		return Statistics{}, err
	}

	if res.Result == "success" {
		var r Statistics
		err := t.extractArgs(res, &r)
		if err != nil {
			return Statistics{}, err
		}
		return r, err
	}

	return Statistics{}, fmt.Errorf("request failed")
}

func (t *Transmission) SessionInfo() (Info, error) {
	res, err := t.makeCall(&Request{
		Method: "session-get",
	})
	if err != nil {
		return Info{}, err
	}

	if res.Result == "success" {
		var r Info
		err := t.extractArgs(res, &r)
		if err != nil {
			return Info{}, err
		}
		return r, nil
	}

	return Info{}, fmt.Errorf("request failed")
}

func (t *Transmission) FreeSpace(path string) error {
	res, err := t.makeCall(&Request{
		Method:    "free-space",
		Arguments: ReqArguments{Path: path},
	})
	if err != nil {
		return err
	}

	if res.Result == "success" {
		return nil
	}

	return fmt.Errorf("request failed")
}
