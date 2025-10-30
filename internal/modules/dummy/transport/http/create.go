package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

// DummyCreateRequest represents request for dummy creation.
// swagger:parameters dummyCreate
type DummyCreateRequest struct {
	Name string `json:"name"`
}

// swagger:route POST /dummy dummyCreate
// Create new dummy
//
//	Create dummy with provided name.
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Responses:
//	  200: DummyResponse
//	  400: ErrorResponse
//	  500: ErrorResponse
func (d dummyHandlers) createTest(writer http.ResponseWriter, request *http.Request) {
	var req DummyCreateRequest

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		render.JSON(writer, request, ErrorResponse{
			Code:    "bad_request",
			Message: fmt.Sprintf("invalid request body: %v", err),
		})
		return
	}

	if req.Name == "" {
		writer.WriteHeader(http.StatusBadRequest)
		render.JSON(writer, request, ErrorResponse{
			Code:    "bad_request",
			Message: "name is required",
		})
		return
	}

	res, err := d.service.Create(request.Context(), req.Name)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		render.JSON(writer, request, ErrorResponse{
			Code:    "internal",
			Message: fmt.Sprintf("create dummy error: %v", err),
		})
		return
	}

	resp := DummyResponse{
		ID:   res.ID,
		Name: res.Name,
	}

	render.JSON(writer, request, resp)
}
