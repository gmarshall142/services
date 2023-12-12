package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gmarshall142/services/api/models"
	"github.com/gmarshall142/services/api/responses"
	"github.com/gmarshall142/services/api/utils/formaterror"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (server *Server) GetAudioFormats(w http.ResponseWriter, r *http.Request) {

	audioFormat := models.AudioFormat{}

	audioFormats, err := audioFormat.FindAllAudioFormats(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, audioFormats)
}

func (server *Server) GetAudios(w http.ResponseWriter, r *http.Request) {

	audio := models.Audio{}

	audios, err := audio.FindAllAudios(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, audios)
}

func (server *Server) CreateAudio(w http.ResponseWriter, r *http.Request) {
	audio, err := getAudioRecord(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	audioCreated, err := audio.SaveAudio(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, audioCreated.ID))
	responses.JSON(w, http.StatusCreated, audioCreated)
}

func (server *Server) GetAudiosByTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	// TODO: verify title is OK
	//uid, err := strconv.ParseUint(vars["title"], 10, 32)
	//if err != nil {
	//	responses.ERROR(w, http.StatusBadRequest, err)
	//	return
	//}
	audio := models.Audio{}
	audios, err := audio.FindAllAudiosByTitle(server.DB, title)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, audios)
}

func (server *Server) GetAudioData(w http.ResponseWriter, r *http.Request) {

	//vars := mux.Vars(r)
	rawUrl, _ := url.Parse(r.RequestURI)
	//params, _ := url.ParseQuery(rawUrl.RawQuery)
	//uid, err := strconv.ParseUint(vars["id"], 10, 32)
	//if err != nil {
	//	responses.ERROR(w, http.StatusBadRequest, err)
	//	return
	//}
	audio := models.Audio{}
	rec, err := audio.FindAudioByDiscogsSearch(rawUrl.RawQuery)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, rec)
}

func (server *Server) UpdateAudio(w http.ResponseWriter, r *http.Request) {
	// get record ID
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	audio, err := getAudioRecord(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedAudio, err := audio.UpdateAudio(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedAudio)
}

func (server *Server) DeleteAudio(w http.ResponseWriter, r *http.Request) {
	// get record ID
	vars := mux.Vars(r)
	audio := models.Audio{}
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	_, err = audio.DeleteAudio(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}

func getAudioRecord(r *http.Request) (models.Audio, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return models.Audio{}, err
	}
	audio := models.Audio{}
	err = json.Unmarshal(body, &audio)
	if err != nil {
		return models.Audio{}, err
	}
	audio.Prepare()
	err = audio.Validate("update")
	if err != nil {
		return models.Audio{}, err
	}
	return audio, nil
}
