package controllers

import (
	"github.com/gmarshall142/services/api/models"
	"github.com/gmarshall142/services/api/responses"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (server *Server) GetBikeRims(w http.ResponseWriter, r *http.Request) {

	bikeRim := models.BikeRim{}

	bikeRims, err := bikeRim.FindAllBikeRims(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, bikeRims)
}

func (server *Server) GetBikeRim(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	bikeRim := models.BikeRim{}
	rec, err := bikeRim.FindBikeRimByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, rec)
}
