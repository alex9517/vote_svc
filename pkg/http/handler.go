//  Created : 2024-Mar-14
// Modified : 2024-Apr-22

package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	endpoint "vote_svc/pkg/endpoint"
	service "vote_svc/pkg/service"

	"github.com/go-kit/kit/ratelimit"
	http1 "github.com/go-kit/kit/transport/http"
)

// This is not exactly what is usually assumed by the abbr DTO (Data Transfer Object). It's
// just an auxiliary struct for the convenient data handling (see decodeUpdateVoteResultsRequest);
type VoteUpdateDTO struct {
	VoteId      int    `json:"vote_id"`
	ContenderId int16  `json:"co_id"`
	UserId      string `json:"user_id"`
}

//////////////////////////////
//
// MAKE GET VOTE DATA HANDLER
//
//////////////////////////////

// makeGetVoteDataHandler creates the handler logic
func makeGetVoteDataHandler(m *http.ServeMux,
	endpoints endpoint.Endpoints, options []http1.ServerOption) {

	m.Handle("GET /votes/{id}", http1.NewServer(
		endpoints.GetVoteDataEndpoint,
		decodeGetVoteDataRequest,
		encodeGetVoteDataResponse,
		options...))
}

// decodeGetVoteDataRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetVoteDataRequest(_ context.Context, r *http.Request) (interface{}, error) {

	// p1 := r.URL.Query().Get("id")	<-- No! This is absolutely wrong!
	// The following approach (old style) probably works, but there is better way.
	// p1 := strings.TrimPrefix(r.URL.Path, "/votes/")

	// The following is possible if your Go is 1.22 or newer (new approach, I guess).
	p1 := r.PathValue("id")
	if p1 == "" {
		return endpoint.GetVoteDataRequest{}, service.ErrBadRequest
	}

	id, err := strconv.Atoi(p1)
	if err != nil {
		return endpoint.GetVoteDataRequest{}, service.ErrBadRequest
	}

	req := endpoint.GetVoteDataRequest{VoteId: id}
	// err = json.NewDecoder(r.Body).Decode(&req)	<-- This is wrong!
	return req, err
}

// encodeGetVoteDataResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer;

// IMPORTANT!
//
//	For the client JS, this JSON response looks like: rec.v0.vote_id, rec.v0.header, ..
//	I've decided to accept this clumsy 'v0' to avoid extra processing and transformations.
func encodeGetVoteDataResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

/////////////////////////////////
//
// MAKE GET VOTE RESULTS HANDLER
//
/////////////////////////////////

func makeGetVoteResultsHandler(m *http.ServeMux,
	endpoints endpoint.Endpoints, options []http1.ServerOption) {

	m.Handle("GET /votes/{id}/results", http1.NewServer(
		endpoints.GetVoteResultsEndpoint,
		decodeGetVoteResultsRequest,
		encodeGetVoteResultsResponse,
		options...))
}

func decodeGetVoteResultsRequest(_ context.Context, r *http.Request) (interface{}, error) {

	p1 := r.PathValue("id")
	if p1 == "" {
		return endpoint.GetVoteResultsRequest{}, service.ErrBadRequest
	}

	id, err := strconv.Atoi(p1)
	if err != nil {
		return endpoint.GetVoteResultsRequest{}, service.ErrBadRequest
	}

	req := endpoint.GetVoteResultsRequest{VoteId: id}
	return req, err
}

// Remember 'v0' (see comments in the prev function);
func encodeGetVoteResultsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

////////////////////////////////////
//
// MAKE UPDATE VOTE RESULTS HANDLER
//
////////////////////////////////////

func makeUpdateVoteResultsHandler(m *http.ServeMux,
	endpoints endpoint.Endpoints, options []http1.ServerOption) {

	m.Handle("PUT /votes", http1.NewServer(
		endpoints.UpdateVoteResultsEndpoint,
		decodeUpdateVoteResultsRequest,
		encodeUpdateVoteResultsResponse,
		options...))
}

func decodeUpdateVoteResultsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := VoteUpdateDTO{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return endpoint.UpdateVoteResultsRequest{}, err
	}

	var decodedReq = endpoint.UpdateVoteResultsRequest{
		VoteId:      req.VoteId,
		ContenderId: req.ContenderId,
		UserId:      req.UserId,
	}

	return decodedReq, nil
}

func encodeUpdateVoteResultsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

///////////////////////////////////
//
// MAKE GET SERVICE STATUS HANDLER
//
///////////////////////////////////

func makeGetServiceStatusHandler(m *http.ServeMux,
	endpoints endpoint.Endpoints, options []http1.ServerOption) {

	m.Handle("GET /health", http1.NewServer(
		endpoints.GetServiceStatusEndpoint,
		decodeGetServiceStatusRequest,
		encodeGetServiceStatusResponse, options...))
}

func decodeGetServiceStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoint.GetServiceStatusRequest{}, nil
}

func encodeGetServiceStatusResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

/////////////////////////
//
// ERROR ENCODER/DECODER
//
/////////////////////////

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func ErrorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

// This is used to set the http status, see an example here :
// https://github.com/go-kit/kit/blob/master/examples/addsvc/pkg/addtransport/http.go#L133
func err2code(err error) int {
	switch err {
	case ratelimit.ErrLimited:
		return http.StatusTooManyRequests

	case service.ErrBadRequest:
		return http.StatusBadRequest

	case service.ErrMethodNotAllowed:
		return http.StatusMethodNotAllowed

	case service.ErrNotFound:
		return http.StatusNotFound

	case service.ErrServiceUnavailable:
		return http.StatusServiceUnavailable

	case service.ErrUnauthorized:
		return http.StatusUnauthorized

	case service.ErrForbidden:
		return http.StatusForbidden

	case service.ErrNoContent:
		return http.StatusNoContent

	default:
		return http.StatusInternalServerError
	}
}

type errorWrapper struct {
	Error string `json:"error"`
}

// --- END OF FILE ---
