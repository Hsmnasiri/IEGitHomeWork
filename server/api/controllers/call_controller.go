package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Hsmnasiri/http_monitoring/server/api/models"
	"github.com/Hsmnasiri/http_monitoring/server/api/utils/formaterror"
	"github.com/Hsmnasiri/http_monitoring/server/api/utils/responses"

	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (server *Server) CreateCall(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	epc := models.EndPointCalls{}
	err = json.Unmarshal(body, &epc)
	fmt.Println(epc)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	epc.Prepare()
	err = epc.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }
	// if uid != epc.OwnerID {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }
	epcCreated, err := epc.SaveCall(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, epcCreated.ID))
	responses.JSON(w, http.StatusCreated, epcCreated)
}

func (server *Server) GetCalls(w http.ResponseWriter, r *http.Request) {

	epc := models.EndPointCalls{}

	epcs, err := epc.FindAllCalls(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, epcs)
}
func (server *Server) GetCallsByTime(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	type timeQuery struct {
		urlID     string
		StartTime time.Time
		EndTime   time.Time
	}
	tq := new(timeQuery)
	err = json.Unmarshal(body, &tq)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// uid, err := strconv.ParseUint(tq.urlID, 10, 32)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }
	epc := models.EndPointCalls{}
	epcs, err := epc.FindCallsByTime(server.DB, tq.StartTime, tq.EndTime)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, epcs)
}

func (server *Server) GetCall(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	epcid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	epc := models.EndPointCalls{}

	CallReceived, err := epc.FindCallByID(server.DB, uint32(epcid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, CallReceived)
}
