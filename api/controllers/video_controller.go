package controllers

import (
	"github.com/gmarshall142/services/api/models"
	"github.com/gmarshall142/services/api/responses"
	"github.com/gorilla/mux"
	"net/http"
)

//func (server *Server) GetBikes(w http.ResponseWriter, r *http.Request) {
//
//	bike := models.Bike{}
//
//	bikes, err := bike.FindAllBikes(server.DB)
//	if err != nil {
//		responses.ERROR(w, http.StatusInternalServerError, err)
//		return
//	}
//	responses.JSON(w, http.StatusOK, bikes)
//}
//
//func (server *Server) CreateBike(w http.ResponseWriter, r *http.Request) {
//	bike, err := getBikeRecord(r)
//	if err != nil {
//		responses.ERROR(w, http.StatusUnprocessableEntity, err)
//		return
//	}
//
//	bikeCreated, err := bike.SaveBike(server.DB)
//	if err != nil {
//		formattedError := formaterror.FormatError(err.Error())
//		responses.ERROR(w, http.StatusInternalServerError, formattedError)
//		return
//	}
//
//	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, bikeCreated.ID))
//	responses.JSON(w, http.StatusCreated, bikeCreated)
//}

func (server *Server) GetVideo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	//uid, err := strconv.ParseUint(vars["id"], 10, 32)
	//if err != nil {
	//	responses.ERROR(w, http.StatusBadRequest, err)
	//	return
	//}
	video := models.Video{}
	rec, err := video.FindVideoByImdbID(server.DB, vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, rec)
}

//func (server *Server) UpdateBike(w http.ResponseWriter, r *http.Request) {
//	// get record ID
//	vars := mux.Vars(r)
//	uid, err := strconv.ParseUint(vars["id"], 10, 32)
//	if err != nil {
//		responses.ERROR(w, http.StatusBadRequest, err)
//		return
//	}
//
//	bike, err := getBikeRecord(r)
//	if err != nil {
//		responses.ERROR(w, http.StatusUnprocessableEntity, err)
//		return
//	}
//
//	updatedBike, err := bike.UpdateBike(server.DB, uint32(uid))
//	if err != nil {
//		formattedError := formaterror.FormatError(err.Error())
//		responses.ERROR(w, http.StatusInternalServerError, formattedError)
//		return
//	}
//	responses.JSON(w, http.StatusOK, updatedBike)
//}
//
//func (server *Server) DeleteBike(w http.ResponseWriter, r *http.Request) {
//	// get record ID
//	vars := mux.Vars(r)
//	bike := models.Bike{}
//	uid, err := strconv.ParseUint(vars["id"], 10, 32)
//	if err != nil {
//		responses.ERROR(w, http.StatusBadRequest, err)
//		return
//	}
//	_, err = bike.DeleteBike(server.DB, uint32(uid))
//	if err != nil {
//		responses.ERROR(w, http.StatusInternalServerError, err)
//		return
//	}
//	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
//	responses.JSON(w, http.StatusNoContent, "")
//}

//func getBikeRecord(r *http.Request) (models.Bike, error) {
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		return models.Bike{}, err
//	}
//	bike := models.Bike{}
//	err = json.Unmarshal(body, &bike)
//	if err != nil {
//		return models.Bike{}, err
//	}
//	bike.Prepare()
//	err = bike.Validate("update")
//	if err != nil {
//		return models.Bike{}, err
//	}
//	return bike, nil
//}
