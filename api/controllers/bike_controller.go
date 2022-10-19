package controllers

import (
	"github.com/gmarshall142/services/api/models"
	"github.com/gmarshall142/services/api/responses"
	"github.com/gorilla/mux"
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
