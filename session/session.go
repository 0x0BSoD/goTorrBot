package session

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/0x0bsod/torrBot/torrent"
	"log"
	"net/http"
	"time"
)

// Session - main struct
type Session struct {
	*Client
}

// NewSession return session instance with token
func NewSession(p Parameters) (*Session, error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := Client{
		Url:      p.Url,
		User:     p.User,
		Password: p.Password,
		Debug:    p.Debug,
		Http:     &http.Client{Transport: tr},
	}

	resp, err := client.getResponse("GET", "/", nil)
	if err != nil {
		return nil, fmt.Errorf("error during getting token: %s", err)
	}

	client.Token = resp.Header.Get("X-Transmission-Session-Id")

	return &Session{&client}, nil
}

func (s *Session) WrappedCall(r *Request) (*Response, error) {
	data, err := s.ApiCall(r)
	if err != nil {
		return nil, err
	}

	if s.Debug {
		log.Printf("[DEBUG] %s", string(data))
	}

	var res Response
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func ExtractArgs(res *Response, result interface{}) error {
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
// =====================================================================================================================

func (s *Session) GetAllTorrents() ([]*torrent.Torrent, error) {
	res, err := s.WrappedCall(&Request{
		Method: "torrent-get",
		Arguments: ReqArguments{
			Fields: torrent.FieldList(
				torrent.ID,
				torrent.Name,
				torrent.Status,
				torrent.Comment,
				torrent.Error,
				torrent.ErrorString,
				torrent.IsFinished,
				torrent.LeftUntilDone,
				torrent.PercentDone,
				torrent.Eta,
				torrent.SizeWhenDone,
				torrent.StartDate,
				torrent.UploadRatio,
				torrent.TotalSize),
		},
	})
	if err != nil {
		return []*torrent.Torrent{}, err
	}

	if res.Result == "success" {
		var r torrent.Torrents
		err := ExtractArgs(res, &r)
		if err != nil {
			return []*torrent.Torrent{}, err
		}
		if len(r.Torrents) == 0 {
			return r.Torrents, fmt.Errorf("no torrents")
		}

		r.ResolveStatus()

		return r.Torrents, nil
	}

	return []*torrent.Torrent{}, fmt.Errorf("request failed")
}
