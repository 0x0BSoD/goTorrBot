package transmissionRPC

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// https://github.com/transmission/transmission/blob/master/extras/rpc-spec.txt

type Transmission struct {
	Http        *Client
	DownloadDir string
	Paused      bool
	Debug       bool
}

func (t *Transmission) makeCall(r *Request) (*Response, error) {
	data, err := t.Http.apiCall(r)
	if err != nil {
		return nil, err
	}

	if t.Debug {
		log.Printf("[DEBUG] %s", string(data))
	}

	var res Response
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (t *Transmission) extractArgs(res *Response, result interface{}) error {
	tmp, err := json.Marshal(res.Arguments)
	if err != nil {
		return err
	}

	fmt.Println("here")

	err = json.Unmarshal(tmp, result)
	if err != nil {
		return err
	}

	return nil
}

// =====================================================================================================================
// Get
// =====================================================================================================================

func (t *Transmission) All() (Torrent, error) {
	res, err := t.makeCall(&Request{
		Method: "torrent-get",
		Arguments: ReqArguments{
			Fields: FieldList(ID, Name, Status, Comment,
				Error, ErrorString, IsFinished, LeftUntilDone,
				PercentDone, Eta, SizeWhenDone, StartDate,
				UploadRatio, TotalSize),
		},
	})
	if err != nil {
		return Torrent{}, err
	}

	if res.Result == "success" {
		var r Torrent
		err := t.extractArgs(res, &r)
		fmt.Println(err)
		if err != nil {
			return Torrent{}, err
		}
		return r, nil
	}

	return Torrent{}, fmt.Errorf("request failed")
}

func (t *Transmission) ByID(ID int) (Torrent, error) {
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
		return Torrent{}, err
	}

	if res.Result == "success" {
		var r Torrent
		err := t.extractArgs(res, &r)
		if err != nil {
			return Torrent{}, err
		}
		return r, nil
	}

	return Torrent{}, fmt.Errorf("request failed")
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
		Format: "table",
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
		Format: "table",
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

func (t *Transmission) Verify(ID int) error {
	res, err := t.makeCall(&Request{
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

func (t *Transmission) Start(ID int) error {
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

func (t *Transmission) Stop(ID int) error {
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

func (t *Transmission) Remove(ID int, rmLocalData bool) error {
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
