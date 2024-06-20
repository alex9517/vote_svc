//  Created : 2024-Mar-14
// Modified : 2024-Apr-24

package endpoint

import (
	"context"
	"strconv"
	"time"
	service "vote_svc/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
	"github.com/patrickmn/go-cache"
)

///////////////////////////////
//
// MAKE GET VOTE DATA ENDPOINT
//
///////////////////////////////

// GetVoteDataRequest collects the request parameters for the GetVoteData method.
type GetVoteDataRequest struct {
	VoteId int `json:"vote_id"`
}

// GetVoteDataResponse collects the response parameters for the GetVoteData method.
type GetVoteDataResponse struct {
	V0 *service.VoteData `json:"v0"`
	E1 error             `json:"e1"`
}

// MakeGetVoteDataEndpoint returns an endpoint that invokes GetVoteData on the service.
func MakeGetVoteDataEndpoint(s service.VoteService, c *cache.Cache) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetVoteDataRequest)

		v0, ok := c.Get(strconv.Itoa(req.VoteId))
		if ok {
			return GetVoteDataResponse{
				V0: v0.(*service.VoteData),
				E1: nil,
			}, nil
		}

		v0, e1 := s.GetVoteData(ctx, req.VoteId)
		c.Set(strconv.Itoa(req.VoteId), v0, cache.DefaultExpiration)
		return GetVoteDataResponse{
			E1: e1,
			V0: v0.(*service.VoteData),
		}, nil
	}
}

// Failed implements Failer.
func (r GetVoteDataResponse) Failed() error {
	return r.E1
}

//////////////////////////////////
//
// MAKE GET VOTE RESULTS ENDPOINT
//
//////////////////////////////////

// GetVoteResultsRequest collects the request parameters for the GetVoteResults method.
type GetVoteResultsRequest struct {
	VoteId int `json:"vote_id"`
}

// GetVoteResultsResponse collects the response parameters for the GetVoteResults method.
type GetVoteResultsResponse struct {
	V0 *service.VoteData `json:"v0"`
	E1 error             `json:"e1"`
}

// MakeGetVoteResultsEndpoint returns an endpoint that invokes GetVoteResults on the service.
func MakeGetVoteResultsEndpoint(s service.VoteService, c *cache.Cache) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetVoteResultsRequest)

		key := strconv.Itoa(req.VoteId) + "r"
		v0, ok := c.Get(key)
		if ok {
			return GetVoteDataResponse{
				V0: v0.(*service.VoteData),
				E1: nil,
			}, nil
		}

		v0, e1 := s.GetVoteData(ctx, req.VoteId)
		c.Set(key, v0, 10*time.Second)
		return GetVoteResultsResponse{
			E1: e1,
			V0: v0.(*service.VoteData),
		}, nil
	}
}

// Failed implements Failer.
func (r GetVoteResultsResponse) Failed() error {
	return r.E1
}

/////////////////////////////////////
//
// MAKE UPDATE VOTE RESULTS ENDPOINT
//
/////////////////////////////////////

// UpdateVoteResultsRequest collects the request parameters for the UpdateVoteResults method.
type UpdateVoteResultsRequest struct {
	VoteId      int    `json:"vote_id"`
	ContenderId int16  `json:"co_id"`
	UserId      string `json:"user_id"`
}

// UpdateVoteResultsResponse collects the response parameters for the UpdateVoteResults method.
type UpdateVoteResultsResponse struct {
	E0 error `json:"e0"`
}

// MakeUpdateVoteResultsEndpoint returns an endpoint that invokes UpdateVoteResults on the service.
func MakeUpdateVoteResultsEndpoint(s service.VoteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateVoteResultsRequest)
		e0 := s.UpdateVoteResults(ctx, req.VoteId, req.ContenderId, req.UserId)
		return UpdateVoteResultsResponse{E0: e0}, nil
	}
}

// Failed implements Failer.
func (r UpdateVoteResultsResponse) Failed() error {
	return r.E0
}

////////////////////////////////////
//
// MAKE GET SERVICE STATUS ENDPOINT
//
////////////////////////////////////

// GetServiceStatusRequest collects the request parameters for the GetServiceStatus method.
type GetServiceStatusRequest struct{}

// GetServiceStatusResponse collects the response parameters for the GetServiceStatus method.
type GetServiceStatusResponse struct {
	H0 *service.HealthStatus `json:"h0"`
}

// MakeGetServiceStatusEndpoint returns an endpoint that invokes GetServiceStatus on the service.
func MakeGetServiceStatusEndpoint(s service.VoteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		h0 := s.GetServiceStatus(ctx)
		return GetServiceStatusResponse{H0: h0}, nil
	}
}

// Failed implements Failer.
func (r GetServiceStatusResponse) Failed() error {
	if r.H0.Status == 200 {
		return nil
	}

	return service.ErrServiceUnavailable
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

/*
// GetVoteData implements Service. Primarily useful in a client.
func (e Endpoints) GetVoteData(ctx context.Context, vote_id int) (v0 *service.VoteData, e1 error) {
	request := GetVoteDataRequest{VoteId: vote_id}
	response, err := e.GetVoteDataEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetVoteDataResponse).V0, response.(GetVoteDataResponse).E1
}

// GetVoteResults implements Service. Primarily useful in a client.
func (e Endpoints) GetVoteResults(ctx context.Context, vote_id int32) (v0 *service.VoteResults, e1 error) {
	request := GetVoteResultsRequest{VoteId: vote_id}
	response, err := e.GetVoteResultsEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetVoteResultsResponse).V0, response.(GetVoteResultsResponse).E1
}

// UpdateVoteResults implements Service. Primarily useful in a client.
func (e Endpoints) UpdateVoteResults(ctx context.Context, vote_id int, co_id int16, user_id string) (e0 error) {
	request := UpdateVoteResultsRequest{
		ContenderId: co_id,
		UserId:      user_id,
		VoteId:      vote_id,
	}
	response, err := e.UpdateVoteResultsEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(UpdateVoteResultsResponse).E0
}

// GetServiceStatus implements Service. Primarily useful in a client.
func (e Endpoints) GetServiceStatus(ctx context.Context) (h0 *service.HealthStatus) {
	request := GetServiceStatusRequest{}
	response, err := e.GetServiceStatusEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetServiceStatusResponse).H0
}
*/

// --- END ---
