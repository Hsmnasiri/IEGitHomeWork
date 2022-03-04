package controllers

import (
	"net/http"

	"github.com/Hsmnasiri/http_monitoring/server/api/utils/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To The HTTP Monitoring  App")

}
