package http

import (
	"net/http"

	"github.com/go-chi/render"
)

type healthResponse struct {
	Status string `json:"status"`
}

func health(writer http.ResponseWriter, request *http.Request) {
	render.JSON(writer, request, healthResponse{Status: "ok"})
}
