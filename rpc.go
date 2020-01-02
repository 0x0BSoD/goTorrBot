package main

import (
	"encoding/json"
	"fmt"
)

// https://github.com/transmission/transmission/blob/master/extras/rpc-spec.txt

type Transmission struct {
	http *Client
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

	fmt.Println("[DEBUG]", string(data))

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

	fmt.Println("[DEBUG]", string(data))

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

	fmt.Println("[DEBUG]", string(data))

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

//func (t *Transmission) AddMagnet(magnetLink string) (Response, error) {
//	//var names []string
//	//var peers []string
//	//
//	//for _, i := range strings.Split(magnetLink, "&") {
//	//	decoded, err := url.QueryUnescape(i)
//	//	if err != nil {
//	//		panic(err)
//	//	}
//	//	if strings.HasPrefix(decoded, "dn=") {
//	//		names = append(names, strings.ReplaceAll(decoded, "dn=", ""))
//	//	}
//	//	if strings.HasPrefix(decoded, "tr=") {
//	//		peers = append(peers, strings.ReplaceAll(decoded, "tr=", ""))
//	//	}
//	//}
//	//
//	//fmt.Println(names)
//	//fmt.Println(peers)
//
//	p := &Request{
//		Method: "torrent-add",
//		Arguments: ReqArguments{
//			DownloadDir: "/home/alex",
//			FileName:    magnetLink,
//			Paused:      true,
//		},
//		Format: "table",
//	}
//
//	data, err := t.http.apiCall(p)
//	if err != nil {
//		return Response{}, err
//	}
//
//	//fmt.Println("[DEBUG]", string(data))
//
//	var result Response
//	err = json.Unmarshal(data, &result)
//	if err != nil {
//		return Response{}, err
//	}
//
//	if result.Arguments.TorrentAdded.ID != 0 {
//		_, err = t.Verify(result.Arguments.TorrentAdded.ID)
//		if err != nil {
//			return Response{}, err
//		}
//	}
//
//
//	return result, nil
//}

func (t *Transmission) Verify(ID int) (Response, error) {
	p := &Request{
		Method: "torrent-verify",
		Arguments: ReqArguments{
			IDs: []int{ID},
		},
	}

	data, err := t.http.apiCall(p)
	if err != nil {
		return Response{}, err
	}

	fmt.Println("[DEBUG]", string(data))

	var result Response
	err = json.Unmarshal(data, &result)
	if err != nil {
		return Response{}, err
	}

	return result, nil
}

func (t *Transmission) Start(ID int) (Response, error) {
	p := &Request{
		Method: "torrent-start",
		Arguments: ReqArguments{
			IDs: []int{ID},
		},
	}

	data, err := t.http.apiCall(p)
	if err != nil {
		return Response{}, err
	}

	fmt.Println("[DEBUG]", string(data))

	var result Response
	err = json.Unmarshal(data, &result)
	if err != nil {
		return Response{}, err
	}

	return result, nil
}

func (t *Transmission) Stop(ID int) (Response, error) {
	p := &Request{
		Method: "torrent-stop",
		Arguments: ReqArguments{
			IDs: []int{ID},
		},
	}

	data, err := t.http.apiCall(p)
	if err != nil {
		return Response{}, err
	}

	fmt.Println("[DEBUG]", string(data))

	var result Response
	err = json.Unmarshal(data, &result)
	if err != nil {
		return Response{}, err
	}

	return result, nil
}
