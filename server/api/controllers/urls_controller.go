package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Hsmnasiri/http_monitoring/server/api/auth"
	"github.com/Hsmnasiri/http_monitoring/server/api/models"
	"github.com/Hsmnasiri/http_monitoring/server/api/utils/formaterror"
	"github.com/Hsmnasiri/http_monitoring/server/api/utils/responses"

	"github.com/gorilla/mux"
)

func (server *Server) Createurl(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	url := models.Urls{}
	err = json.Unmarshal(body, &url)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	url.Prepare()
	err = url.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != url.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	urlCreated, err := url.SaveUrl(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, urlCreated.ID))
	responses.JSON(w, http.StatusCreated, urlCreated)
}

func (server *Server) Geturls(w http.ResponseWriter, r *http.Request) {

	url := models.Urls{}

	urls, err := url.FindAllUrlses(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, urls)
}

func (server *Server) Geturl(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	url := models.Urls{}

	urlReceived, err := url.FindUrlByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, urlReceived)
}

func (server *Server) Updateurl(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the url id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the url exist
	url := models.Urls{}
	err = server.DB.Debug().Model(models.Urls{}).Where("id = ?", pid).Take(&url).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("url not found"))
		return
	}

	if uid != url.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	urlUpdate := models.Urls{}
	err = json.Unmarshal(body, &urlUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if uid != urlUpdate.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	urlUpdate.Prepare()
	err = urlUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	urlUpdate.ID = url.ID

	urlUpdated, err := urlUpdate.UpdateAUrl(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, urlUpdated)
}

func (server *Server) Deleteurl(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	url := models.Urls{}
	err = server.DB.Debug().Model(models.Urls{}).Where("id = ?", pid).Take(&url).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	if uid != url.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = url.DeleteAUrl(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}

// func (server *Server) IncrementFailed(url *models.Urls) error {
// 	url.FailedTimes += 1
// 	return server.Updateurl(url)
// }

// func (server *Server) IncrementSuccess(url *models.Urls) error {
// 	url.SuccessTimes += 1
// 	return server.UpdateUrl(url)
// }

// func (server *Server) AddRequest(req *models.EndPointCalls) error {
// 	return server.db.Create(req).Error
// }

// func (server *Server) GetEndPointCallesByUrl(uid uint32) ([]models.EndPointCalls, error) {
// 	var EndPointCalles []models.EndPointCalls
// 	if err := server.db.Model(&models.EndPointCalls{urlID: uid}).Where("url_id == ?", uid).Find(&EndPointCalles).Error; err != nil {
// 		return nil, err
// 	}
// 	return EndPointCalles, nil
// }
