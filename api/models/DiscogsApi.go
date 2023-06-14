package models

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type RawAudio struct {
	MasterID    uint
	Title       string
	ImageUrl    string
	Genres      pq.StringArray
	Artists     pq.StringArray
	Catno       string
	Barcode     string
	Year        string
	AudioTracks []AudioTrack
}

type AudioBaseInfo struct {
	MasterID uint     `json:"master_id"`
	Title    string   `json:"title"`
	Catno    string   `json:"catno"`
	Barcode  []string `json:"barcode"`
	Thumb    string   `json:"thumb"`
	Year     string   `json:"year"`
}

type AudioBaseInfoResults struct {
	Results []AudioBaseInfo `json:"results"`
}

type AudioArtist struct {
	Name string `json:"name"`
}
type AudioTrackList struct {
	Title    string `json:"title"`
	Position string `json:"position"`
	Duration string `json:"duration"`
}
type AudioMasterResults struct {
	Title     string           `json:"title"`
	Artists   []AudioArtist    `json:"artists"`
	Genres    []string         `json:"genres"`
	Styles    []string         `json:"styles"`
	TrackList []AudioTrackList `json:"tracklist"`
}

func getDiscogsRecord(params url.Values) (*RawAudio, error) {
	fmt.Println(params)
	catno := params.Get("catno")
	fmt.Println(catno)
	barcode := params.Get("barcode")
	fmt.Println(barcode)

	audio := RawAudio{}

	var url string
	// Search
	if catno != "" {
		url = "https://api.discogs.com/database/search?catno=" + catno
	} else {
		url = "https://api.discogs.com/database/search?barcode=" + barcode
	}
	bodyBytes, err := discogsCall(url)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	var responseObject AudioBaseInfoResults
	json.Unmarshal(bodyBytes, &responseObject)

	if len(responseObject.Results) > 0 {
		results := responseObject.Results[0]
		audio.MasterID = results.MasterID
		audio.Catno = results.Catno
		if len(results.Barcode) > 0 {
			audio.Barcode = results.Barcode[0]
		}
		audio.Year = results.Year
		audio.ImageUrl = results.Thumb
		// Master Query
		url = "https://api.discogs.com/masters/" + strconv.FormatUint(uint64(audio.MasterID), 10)
		bodyBytes, err = discogsCall(url)
		if err != nil {
			fmt.Print(err.Error())
			return nil, err
		}
		var masterObj AudioMasterResults
		json.Unmarshal(bodyBytes, &masterObj)
		audio.Title = masterObj.Title
		for _, artist := range masterObj.Artists {
			audio.Artists = append(audio.Artists, artist.Name)
		}
		for _, genre := range masterObj.Genres {
			audio.Genres = append(audio.Genres, genre)
		}
		for _, style := range masterObj.Styles {
			audio.Genres = append(audio.Genres, style)
		}
		for _, track := range masterObj.TrackList {
			audioTrack := AudioTrack{}
			audioTrack.Position = track.Position
			audioTrack.Title = track.Title
			arr := strings.Split(track.Duration, ":")
			if len(arr) == 2 {
				min, minErr := strconv.ParseInt(arr[0], 0, 32)
				sec, secErr := strconv.ParseInt(arr[1], 0, 32)
				if minErr == nil && secErr == nil {
					audioTrack.Duration = uint(min*60 + sec)
				}
			}
			audio.AudioTracks = append(audio.AudioTracks, audioTrack)
		}
	}

	return &audio, nil
}

func discogsCall(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	key := os.Getenv("DISCOGS_KEY")
	secret := os.Getenv("DISCOGS_SECRET")
	discogsStr := "Discogs key=" + key + ", secret=" + secret
	req.Header.Add("Authorization", discogsStr)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	return bodyBytes, nil
}
