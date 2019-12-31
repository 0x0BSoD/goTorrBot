package main

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Log("Creating NewClient")
	c, err := NewClient(srvAddr, user, pass)
	if err != nil {
		t.Error(err)
	}

	if c.Token == "" {
		t.Error("token must be set")
	}

	fmt.Printf("%+v\n", c)

	r, err := c.SessionStats()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", r)
}

func TestClient_SessionStats(t *testing.T) {
	c, err := NewClient(srvAddr, user, pass)
	if err != nil {
		t.Error(err)
	}
	r, err := c.SessionStats()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", r)
}

func TestClient_TorrentGet(t *testing.T) {
	c, err := NewClient(srvAddr, user, pass)
	if err != nil {
		t.Error(err)
	}

	c.TorrentGet()
}