package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type RawVideo struct {
	ImdbID      string
	Name        string
	ImageUrl    string
	ImageWidth  uint
	ImageHeight uint
	Runtime     uint
	Genres      []string
	Plot        string
	Actors      []string
}
type TitleText struct {
	Text     string `json:"text"`
	TypeName string `json:"__typename"`
}

type PrimaryImage struct {
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
	Url    string `json:"url"`
}

type Genre struct {
	Text string `json:"text"`
}

type Genres struct {
	Genres []Genre `json:"genres"`
}

type Runtime struct {
	Seconds uint `json:"seconds"`
}

type PlotText struct {
	Text string `json:"plainText"`
}
type Plot struct {
	PlotText PlotText `json:"plotText"`
}

type Name struct {
	NameText TitleText `json:"nameText"`
}
type Credits struct {
	Name Name `json:"name"`
}
type PrincipalCast struct {
	Credits []Credits `json:"credits"`
}

type BaseInfo struct {
	ID            string          `json:"id"`
	TitleText     TitleText       `json:"titleText"`
	Image         PrimaryImage    `json:"primaryImage"`
	Runtime       Runtime         `json:"runtime"`
	Genres        Genres          `json:"genres"`
	Plot          Plot            `json:"plot"`
	PrincipalCast []PrincipalCast `json:"principalCast"`
}

type BaseInfoResults struct {
	Results BaseInfo `json:"results"`
}

func getMoviesDbRecord(id string) (*RawVideo, error) {
	video := RawVideo{}

	//============================================================================
	baseObj, err := moviesDbRequest(id, "base_info")
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	video.ImdbID = baseObj.Results.ID
	video.Name = baseObj.Results.TitleText.Text
	video.ImageUrl = baseObj.Results.Image.Url
	video.ImageWidth = baseObj.Results.Image.Width
	video.ImageHeight = baseObj.Results.Image.Height
	video.Runtime = baseObj.Results.Runtime.Seconds
	for _, genre := range baseObj.Results.Genres.Genres {
		video.Genres = append(video.Genres, genre.Text)
	}
	video.Plot = baseObj.Results.Plot.PlotText.Text

	//============================================================================
	respObj, err := moviesDbRequest(id, "principalCast")
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	for _, pc := range respObj.Results.PrincipalCast {
		for _, credit := range pc.Credits {
			video.Actors = append(video.Actors, credit.Name.NameText.Text)
		}
	}

	return &video, nil
}

func moviesDbRequest(id string, info string) (*BaseInfoResults, error) {
	url := "https://moviesdatabase.p.rapidapi.com/titles/" + id + "?info=" + info
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-RapidAPI-Key", os.Getenv("RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", "moviesdatabase.p.rapidapi.com")
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

	var responseObject BaseInfoResults
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response as struct %+v\n", responseObject)

	return &responseObject, nil
}
