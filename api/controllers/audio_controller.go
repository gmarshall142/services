package controllers

import (
	"github.com/gmarshall142/services/api/models"
	"github.com/gmarshall142/services/api/responses"
	"net/http"
	"net/url"
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

//func (server *Server) GetVideos(w http.ResponseWriter, r *http.Request) {
//
//	video := models.Video{}
//
//	videos, err := video.FindAllVideos(server.DB)
//	if err != nil {
//		responses.ERROR(w, http.StatusInternalServerError, err)
//		return
//	}
//	responses.JSON(w, http.StatusOK, videos)
//}
//
//func (server *Server) CreateVideo(w http.ResponseWriter, r *http.Request) {
//	video, err := getVideoRecord(r)
//	if err != nil {
//		responses.ERROR(w, http.StatusUnprocessableEntity, err)
//		return
//	}
//
//	videoCreated, err := video.SaveVideo(server.DB)
//	if err != nil {
//		formattedError := formaterror.FormatError(err.Error())
//		responses.ERROR(w, http.StatusInternalServerError, formattedError)
//		return
//	}
//
//	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, videoCreated.ID))
//	responses.JSON(w, http.StatusCreated, videoCreated)
//}
//
//func (server *Server) GetVideosByTitle(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	title := vars["title"]
//	// TODO: verify title is OK
//	//uid, err := strconv.ParseUint(vars["title"], 10, 32)
//	//if err != nil {
//	//	responses.ERROR(w, http.StatusBadRequest, err)
//	//	return
//	//}
//	video := models.Video{}
//	videos, err := video.FindAllVideosByTitle(server.DB, title)
//	if err != nil {
//		responses.ERROR(w, http.StatusInternalServerError, err)
//		return
//	}
//	responses.JSON(w, http.StatusOK, videos)
//}

func (server *Server) GetAudioData(w http.ResponseWriter, r *http.Request) {

	//vars := mux.Vars(r)
	rawUrl, _ := url.Parse(r.RequestURI)
	params, _ := url.ParseQuery(rawUrl.RawQuery)
	//uid, err := strconv.ParseUint(vars["id"], 10, 32)
	//if err != nil {
	//	responses.ERROR(w, http.StatusBadRequest, err)
	//	return
	//}
	audio := models.Audio{}
	rec, err := audio.FindAudioByDiscogsSearch(params)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, rec)
}

//func (server *Server) UpdateVideo(w http.ResponseWriter, r *http.Request) {
//	// get record ID
//	vars := mux.Vars(r)
//	uid, err := strconv.ParseUint(vars["id"], 10, 32)
//	if err != nil {
//		responses.ERROR(w, http.StatusBadRequest, err)
//		return
//	}
//
//	video, err := getVideoRecord(r)
//	if err != nil {
//		responses.ERROR(w, http.StatusUnprocessableEntity, err)
//		return
//	}
//
//	updatedVideo, err := video.UpdateVideo(server.DB, uint32(uid))
//	if err != nil {
//		formattedError := formaterror.FormatError(err.Error())
//		responses.ERROR(w, http.StatusInternalServerError, formattedError)
//		return
//	}
//	responses.JSON(w, http.StatusOK, updatedVideo)
//}
//
//func (server *Server) DeleteVideo(w http.ResponseWriter, r *http.Request) {
//	// get record ID
//	vars := mux.Vars(r)
//	video := models.Video{}
//	uid, err := strconv.ParseUint(vars["id"], 10, 32)
//	if err != nil {
//		responses.ERROR(w, http.StatusBadRequest, err)
//		return
//	}
//	_, err = video.DeleteVideo(server.DB, uint32(uid))
//	if err != nil {
//		responses.ERROR(w, http.StatusInternalServerError, err)
//		return
//	}
//	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
//	responses.JSON(w, http.StatusNoContent, "")
//}
//
//func getVideoRecord(r *http.Request) (models.Video, error) {
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		return models.Video{}, err
//	}
//	video := models.Video{}
//	err = json.Unmarshal(body, &video)
//	if err != nil {
//		return models.Video{}, err
//	}
//	video.Prepare()
//	err = video.Validate("update")
//	if err != nil {
//		return models.Video{}, err
//	}
//	return video, nil
//}
