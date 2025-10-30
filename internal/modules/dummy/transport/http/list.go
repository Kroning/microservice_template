package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"{{index .App "git"}}/internal/modules/dummy"
)

const (
	paramID     = "id"
	paramName   = "name"
	paramOffset = "offset"
	paramLimit  = "limit"
)

// DummyListResponse represents response for dummy list request.
// swagger:model DummyListResponse
type DummyListResponse []DummyResponse

// swagger:route GET /dummy dummyList
// Search dummy.
//
//	Find dummy objects by id and name with pagination.
//	Parameters:
//	+ name: id
//	  in: query
//	  description: filter dummy objects by dummy_id
//	  required: false
//	  type: integer
//	+ name: name
//	  in: query
//	  description: filter dummy objects by name
//	  required: false
//	  type: string
//	+ name: offset
//	  in: query
//	  description: page number
//	  required: false
//	  type: integer
//	+ name: limit
//	  in: query
//	  description: found dummy objects
//	  required: false
//	  type: integer
//
//	Produces:
//	- application/json
//
//	Responses:
//	  200: DummyListResponse
//	  400: ErrorResponse
//	  500: ErrorResponse
func (d dummyHandlers) listTest(writer http.ResponseWriter, request *http.Request) {
	req := dummy.ListRequest{
		Offset: 0,
		Limit:  10,
	}

	vals := request.URL.Query()
	if id, err := strconv.Atoi(vals.Get(paramID)); err == nil && id > 0 {
		req.FilterID = int64(id)
	}

	if name := vals.Get(paramName); name != "" {
		req.FilterName = name
	}

	if offset, err := strconv.ParseInt(vals.Get(paramOffset), 10, 32); err == nil && offset > 0 {
		req.Offset = int32(offset)
	}

	if limit, err := strconv.ParseInt(vals.Get(paramLimit), 10, 32); err == nil && limit > 0 {
		req.Limit = int32(limit)
	}

	res, err := d.service.List(request.Context(), req)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		render.JSON(writer, request, ErrorResponse{
			Code:    "internal",
			Message: fmt.Sprintf("list dummies error: %v", err),
		})

		return
	}

	resp := make(DummyListResponse, 0, len(res))
	for i := range res {
		resp = append(resp, DummyResponse{
			ID:   res[i].ID,
			Name: res[i].Name,
		})
	}

	render.JSON(writer, request, resp)
}
