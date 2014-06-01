package main

import (
	"encoding/json"
	"net/http"
)

const (
	TWITCH_BASE_URL = "https://api.twitch.tv/kraken"
)

type StreamSourceTwitch struct {
	BaseUrl string
}

func NewStreamSourceTwitch() *StreamSourceTwitch {
	return &StreamSourceTwitch{TWITCH_BASE_URL}
}

func (sst *StreamSourceTwitch) FetchStream(identifier string) (*StreamTwitch, error) {
	streamTwitch := new(StreamTwitch)
	r, err := http.Get(sst.BaseUrl + "/streams/" + identifier)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.Decode(streamTwitch)

	return streamTwitch, nil
}
