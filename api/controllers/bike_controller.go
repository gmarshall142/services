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
	"strconv"
)

func (server *Server) GetBikes(w http.ResponseWriter, r *http.Request) {

	bike := models.Bike{}

	bikes, err := bike.FindAllBikes(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, bikes)
}

func (server *Server) CreateBike(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	bike := models.Bike{}
	err = json.Unmarshal(body, &bike)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	bike.Prepare()
	err = bike.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	bikeCreated, err := bike.SaveBike(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, bikeCreated.ID))
	responses.JSON(w, http.StatusCreated, bikeCreated)
}

func (server *Server) GetBike(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	bike := models.Bike{}
	rec, err := bike.FindBikeByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, rec)
}
