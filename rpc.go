package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// https://github.com/transmission/transmission/blob/master/extras/rpc-spec.txt

type Transmission struct {
	http        *Client
	DownloadDir string
	Paused      bool
}

func ReturnArguments(res Response, result interface{}) error {
	tmp, err := json.Marshal(res.Arguments)
	if err != nil {
		return err
	}

	err = json.Unmarshal(tmp, result)
	if err != nil {
		return err
	}

	return nil
}

// =====================================================================================================================
// Get
// =====================================================================================================================

func (t *Transmission) SessionStats() (Statistics, error) {
	p := &Request{
		Method: "session-stats",
	}

	data, err := t.http.apiCall(p)
	if err != nil {
		return Statistics{}, err
	}

	//fmt.Println("[DEBUG]", string(data))

	var res Response
	err = json.Unmarshal(data, &res)
	if err != nil {
		return Statistics{}, err
	}

	if res.Result == "success" {
		var r Statistics
		err := ReturnArguments(res, &r)
		if err != nil {
			return Statistics{}, err
		}
		return r, err
	}

	return Statistics{}, fmt.Errorf("request failed")
}

func (t *Transmission) All() (Torrent, error) {
	p := &Request{
		Method: "torrent-get",
		Arguments: ReqArguments{
			Fields: FieldList(ID, Name, Status, Comment,
				Error, ErrorString, IsFinished, LeftUntilDone,
				PercentDone, Eta, SizeWhenDone, StartDate,
				UploadRatio, TotalSize),
		},
		Format: "table",
	}

	data, err := t.http.apiCall(p)
	if err != nil {
		return Torrent{}, err
	}

	//fmt.Println("[DEBUG]", string(data))

	var res Response
	err = json.Unmarshal(data, &res)
	if err != nil {
		return Torrent{}, err
	}

	if res.Result == "success" {
		var r Torrent
		err := ReturnArguments(res, &r)
		if err != nil {
			return Torrent{}, err
		}
		return r, nil
	}

	return Torrent{}, fmt.Errorf("request failed")
}

func (t *Transmission) ByID(ID int) (Torrent, error) {
	p := &Request{
		Method: "torrent-get",
		Arguments: ReqArguments{
			Fields: FieldList(Name, Status, ErrorString, AddedDate, Peers,
				TotalSize, RateDownload, RateUpload, UploadRatio, Files),
			IDs: []int{ID},
		},
	}

	data, err := t.http.apiCall(p)
	if err != nil {
		return Torrent{}, err
	}

	//fmt.Println("[DEBUG]", string(data))

	var res Response
	err = json.Unmarshal(data, &res)
	if err != nil {
		return Torrent{}, err
	}

	if res.Result == "success" {
		var r Torrent
		err := ReturnArguments(res, &r)
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
	f, err := os.Open(path)
	if err != nil {
		return Added{}, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	base64Str := base64.StdEncoding.EncodeToString(content)

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

	data, err := t.http.apiCall(p)
	if err != nil {
		return Added{}, err
	}

	//fmt.Println("[DEBUG]", string(data))

	var res Response
	err = json.Unmarshal(data, &res)
	if err != nil {
		return Added{}, err
	}

	if res.Result == "success" {
		var r Added
		err := ReturnArguments(res, &r)
		if err != nil {
			return Added{}, err
		}

		if (Added{}) == r {
			var d Duplicate
			err := ReturnArguments(res, &d)
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

	data, err := t.http.apiCall(p)
	if err != nil {
		return Added{}, err
	}

	//fmt.Println("[DEBUG]", string(data))

	var res Response
	err = json.Unmarshal(data, &res)
	if err != nil {
		return Added{}, err
	}

	if res.Result == "success" {
		var r Added
		err := ReturnArguments(res, &r)
		if err != nil {
			return Added{}, err
		}

		if (Added{}) == r {
			var d Duplicate
			err := ReturnArguments(res, &d)
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

func (t *Transmission) Verify(ID int) error {
	p := &Request{
		Method: "torrent-verify",
		Arguments: ReqArguments{
			IDs: []int{ID},
		},
	}

	data, err := t.http.apiCall(p)
	if err != nil {
		return err
	}

	//fmt.Println("[DEBUG]", string(data))

	var result Response
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	if result.Result == "success" {
		return nil
	}

	return fmt.Errorf("request failed")
}

func (t *Transmission) Start(ID int) error {
	p := &Request{
		Method: "torrent-start",
		Arguments: ReqArguments{
			IDs: []int{ID},
		},
	}

	data, err := t.http.apiCall(p)
	if err != nil {
		return err
	}

	//fmt.Println("[DEBUG]", string(data))

	var result Response
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	if result.Result == "success" {
		return nil
	}

	return fmt.Errorf("request failed")
}

func (t *Transmission) Stop(ID int) error {
	p := &Request{
		Method: "torrent-stop",
		Arguments: ReqArguments{
			IDs: []int{ID},
		},
	}

	data, err := t.http.apiCall(p)
	if err != nil {
		return err
	}

	//fmt.Println("[DEBUG]", string(data))

	var result Response
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	if result.Result == "success" {
		return nil
	}

	return fmt.Errorf("request failed")
}

func (t *Transmission) Remove(ID int, rmLocalData bool) error {
	p := &Request{
		Method: "torrent-remove",
		Arguments: ReqArguments{
			IDs:             []int{ID},
			DeleteLocalData: rmLocalData,
		},
	}

	data, err := t.http.apiCall(p)
	if err != nil {
		return err
	}

	//fmt.Println("[DEBUG]", string(data))

	var result Response
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	if result.Result == "success" {
		return nil
	}

	return fmt.Errorf("request failed")
}
