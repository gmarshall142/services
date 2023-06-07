package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	MasterID uint
	Title    string
	Year     string
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
	Year     string `json:"year"`
}

// Image         PrimaryImage    `json:"primaryImage"`
// Runtime       Runtime         `json:"runtime"`
// Genres        Genres          `json:"genres"`
// Plot          Plot            `json:"plot"`
// PrincipalCast []PrincipalCast `json:"principalCast"`
// Directors     []Directors     `json:"directors"`
type AudioBaseInfoResults struct {
	Results []AudioBaseInfo `json:"results"`
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
	url := ""
	if catno != "" {
		url = "https://api.discogs.com/database/search?catno=" + catno
	} else {
		url = "https://api.discogs.com/database/search?barcode=" + barcode
	}
	baseObj, err := discogsSearch(url)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	if len(baseObj.Results) > 0 {
		results := baseObj.Results[0]
		audio.MasterID = results.MasterID
		audio.Title = results.Title
		audio.Year = results.Year
		//	video.ImageUrl = baseObj.Results.Image.Url
		//	video.ImageWidth = baseObj.Results.Image.Width
		//	video.ImageHeight = baseObj.Results.Image.Height
		//	video.Runtime = baseObj.Results.Runtime.Seconds
		//	for _, genre := range baseObj.Results.Genres.Genres {
		//		video.Genres = append(video.Genres, genre.Text)
		//	}
		//	video.Plot = baseObj.Results.Plot.PlotText.Text
	}
	//	wg.Done()
	//}(idx)

	// principalCast
	//idx++
	//wg.Add(1)
	//go func(idx int) {
	//	respObj, err := moviesDbRequest(id, "principalCast")
	//	if err != nil {
	//		fmt.Print(err.Error())
	//		return
	//	}
	//	for _, pc := range respObj.Results.PrincipalCast {
	//		for _, credit := range pc.Credits {
	//			video.Actors = append(video.Actors, credit.Name.NameText.Text)
	//		}
	//	}
	//	wg.Done()
	//}(idx)

	// creators_directors_writers
	//idx++
	//wg.Add(1)
	//go func(idx int) {
	//	respObj, err := moviesDbRequest(id, "creators_directors_writers")
	//	if err != nil {
	//		fmt.Print(err.Error())
	//		return
	//	}
	//	for _, dir := range respObj.Results.Directors {
	//		for _, credit := range dir.Credits {
	//			video.Directors = append(video.Directors, credit.Name.NameText.Text)
	//		}
	//	}
	//	wg.Done()
	//}(idx)
	//wg.Wait()

	return &audio, nil
}

func discogsSearch(url string) (*AudioBaseInfoResults, error) {
	//url := "https://api.discogs.com/database/search?" + param + "=" + val
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

	var responseObject AudioBaseInfoResults
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response as struct %+v\n", responseObject)

	return &responseObject, nil
}
