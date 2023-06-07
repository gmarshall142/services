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
)

//	type RawVideo struct {
//		ImdbID      string
//		Name        string
//		ImageUrl    string
//		ImageWidth  uint
//		ImageHeight uint
//		Runtime     uint
//		Genres      []string
//		Plot        string
//		Actors      []string
//		Directors   []string
//	}
type RawAudio struct {
	MasterID    uint
	Title       string
	ImageUrl    string
	ImageWidth  uint
	ImageHeight uint
	Genres      pq.StringArray
	Artists     pq.StringArray
	Catno       string
	Barcode     string
	Year        string
}

//	type TitleText struct {
//		Text     string `json:"text"`
//		TypeName string `json:"__typename"`
//	}
//
//	type PrimaryImage struct {
//		Width  uint   `json:"width"`
//		Height uint   `json:"height"`
//		Url    string `json:"url"`
//	}
//
//	type Genre struct {
//		Text string `json:"text"`
//	}
//
//	type Genres struct {
//		Genres []Genre `json:"genres"`
//	}
//
//	type Runtime struct {
//		Seconds uint `json:"seconds"`
//	}
//
//	type PlotText struct {
//		Text string `json:"plainText"`
//	}
//
//	type Plot struct {
//		PlotText PlotText `json:"plotText"`
//	}
//
//	type Name struct {
//		NameText TitleText `json:"nameText"`
//	}
//
//	type Credits struct {
//		Name Name `json:"name"`
//	}
//
//	type PrincipalCast struct {
//		Credits []Credits `json:"credits"`
//	}
//
//	type Directors struct {
//		Credits []Credits `json:"credits"`
//	}
type AudioBaseInfo struct {
	MasterID uint   `json:"master_id"`
	Title    string `json:"title"`
	Catno    string `json:"catno"`
	Barcode  string `json:"barcode"`
	Thumb    string `json:"thumb"`
	Year     string `json:"year"`
}

type AudioBaseInfoResults struct {
	Results []AudioBaseInfo `json:"results"`
}

type AudioArtist struct {
	Name string `json:"name"`
}
type AudioMasterResults struct {
	Title   string        `json:"title"`
	Artists []AudioArtist `json:"artists"`
	Genres  []string      `json:"genres"`
	Styles  []string      `json:"styles"`
}

func getDiscogsRecord(params url.Values) (*RawAudio, error) {
	fmt.Println(params)
	catno := params.Get("catno")
	fmt.Println(catno)
	barcode := params.Get("barcode")
	fmt.Println(barcode)

	audio := RawAudio{}
	//wg := sync.WaitGroup{}

	// base_info
	//idx := 0
	//wg.Add(1)
	//go func(idx int) {
	var url string
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
		audio.Barcode = results.Barcode
		audio.Year = results.Year
		audio.ImageUrl = results.Thumb
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
