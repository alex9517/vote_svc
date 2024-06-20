//  Created : 2024-Apr-01
// Modified : 2024-Apr-04

// NOTE! I know these test functions are to large. This is not good,
// this is inconvenient. But refactoring takes time. So, maybe next time..
// Anyway, this is a bad example, this is .. antipattern. Though it works.

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	endpoint "vote_svc/pkg/endpoint"
	"vote_svc/pkg/service"

	endpoint1 "github.com/go-kit/kit/endpoint"
	http1 "github.com/go-kit/kit/transport/http"
)

const BAD_USER_ID = "5287eb7d-4813-49e7-8e2d-c8a6ce8b3c4c"
const GOOD_USER_ID = "8e59e85a-66d1-45f6-9816-0560410b61ca"

/////////////
//
// M O C K S
//
/////////////

func makeGetVoteDataEndpointMock() endpoint1.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(endpoint.GetVoteDataRequest)
		if req.VoteId == 1 {
			return endpoint.GetVoteDataResponse{
				V0: &service.VoteData{VoteId: 1, Header: "Test"},
				E1: nil,
			}, nil
		} else if req.VoteId == 2 {
			return endpoint.GetVoteDataResponse{
				V0: &service.VoteData{},
				E1: service.ErrNotFound,
			}, nil
		} else {
			return endpoint.GetVoteDataResponse{
				V0: &service.VoteData{},
				E1: service.ErrInternalServerError,
			}, nil
		}
	}
}

func makeGetVoteResultsEndpointMock() endpoint1.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(endpoint.GetVoteResultsRequest)
		if req.VoteId == 1 {
			return endpoint.GetVoteResultsResponse{
				V0: &service.VoteData{VoteId: 1, Header: "Test"},
				E1: nil,
			}, nil
		} else if req.VoteId == 2 {
			return endpoint.GetVoteResultsResponse{
				V0: &service.VoteData{},
				E1: service.ErrNotFound,
			}, nil
		} else {
			return endpoint.GetVoteResultsResponse{
				V0: &service.VoteData{},
				E1: service.ErrInternalServerError,
			}, nil
		}
	}
}

func makeUpdateVoteResultsEndpointMock() endpoint1.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(endpoint.UpdateVoteResultsRequest)
		if req.VoteId == 1 && req.UserId != BAD_USER_ID {
			return endpoint.UpdateVoteResultsResponse{
				E0: nil,
			}, nil
		} else if req.VoteId == 1 && req.UserId == BAD_USER_ID {
			return endpoint.UpdateVoteResultsResponse{
				E0: service.ErrForbidden,
			}, nil
		} else if req.VoteId == 2 {
			return endpoint.UpdateVoteResultsResponse{
				E0: service.ErrBadRequest,
			}, nil
		} else {
			return endpoint.UpdateVoteResultsResponse{
				E0: service.ErrInternalServerError,
			}, nil
		}
	}
}

var ServiceStatus int

func makeGetServiceStatusEndpointMock() endpoint1.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		// fmt.Printf("ServiceStatus = %d", ServiceStatus)
		if ServiceStatus == 0 {
			return endpoint.GetServiceStatusResponse{
				H0: &service.HealthStatus{
					Status:  http.StatusOK,
					Message: service.SERVICE_STATUS_OK,
				},
			}, nil
		} else if ServiceStatus == 1 {
			return endpoint.GetServiceStatusResponse{
				H0: &service.HealthStatus{
					Status:  http.StatusServiceUnavailable,
					Message: service.SERVICE_STATUS_LOW_MEMORY,
				},
			}, nil
		} else {
			return endpoint.GetServiceStatusResponse{
				H0: &service.HealthStatus{
					Status:  http.StatusServiceUnavailable,
					Message: service.SERVICE_STATUS_NO_DATABASE,
				},
			}, nil
		}
	}
}

func getEndpoints() endpoint.Endpoints {
	eps := endpoint.Endpoints{
		GetVoteDataEndpoint:       makeGetVoteDataEndpointMock(),
		GetVoteResultsEndpoint:    makeGetVoteResultsEndpointMock(),
		UpdateVoteResultsEndpoint: makeUpdateVoteResultsEndpointMock(),
		GetServiceStatusEndpoint:  makeGetServiceStatusEndpointMock(),
	}

	return eps
}

/////////////
//
// T E S T S
//
/////////////

/////////////////////////////////////
//
// TEST HTTP TRANSPORT GET VOTE DATA
//
/////////////////////////////////////

func TestHttpTransportGetVoteData(t *testing.T) {
	testinfo := "test # 1: GetVoteData"
	eps := getEndpoints()
	m := http.NewServeMux()
	makeGetVoteDataHandler(m, eps, []http1.ServerOption{})

	t.Run(testinfo, func(t *testing.T) {
		// Case 1: everything is good, response is http.StatusOK and full struct VodeData;
		u := "/votes/1"
		req := httptest.NewRequest(http.MethodGet, u, nil)

		w := httptest.NewRecorder()
		if w == nil || req == nil {
			t.Fatal("w (or r) is nil")
		}

		m.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s (case # 1) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodGet, u, http.StatusOK, resp.StatusCode)
		}

		ep := endpoint.GetVoteDataResponse{}
		err := json.NewDecoder(resp.Body).Decode(&ep)
		if err != nil {
			t.Errorf("%s, %v", testinfo, err)
		}

		if ep.V0.Header != "Test" {
			t.Errorf("%s (case # 1) failed, expected header \"Test\" but was %s ", testinfo, ep.V0.Header)
		}

		// Case 2: bad vote_id = 2, response must be http.StatusNotFound;
		u = "/votes/2"
		req = httptest.NewRequest(http.MethodGet, u, nil)
		w = httptest.NewRecorder()

		m.ServeHTTP(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("%s (case # 2) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodGet, u, http.StatusNotFound, resp.StatusCode)
		}

		// Case 3: bad vote_id = 3, response must be http.StatusInternalServerError;
		u = "/votes/3"
		req = httptest.NewRequest(http.MethodGet, u, nil)
		w = httptest.NewRecorder()

		m.ServeHTTP(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("%s (case # 3) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodGet, u, http.StatusInternalServerError, resp.StatusCode)
		}

		// Case 4: bad method PUT instead of GET, but vote_id is good, response must be http.StatusMethodNotAllowed;
		u = "/votes/1"
		req = httptest.NewRequest(http.MethodPut, u, nil)
		w = httptest.NewRecorder()

		m.ServeHTTP(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("%s (case # 4) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodPut, u, http.StatusMethodNotAllowed, resp.StatusCode)
		}
	})
}

////////////////////////////////////////
//
// TEST HTTP TRANSPORT GET VOTE RESULTS
//
////////////////////////////////////////

func TestHttpTransportGetVoteResults(t *testing.T) {
	testinfo := "test # 2: GetVoteResults"
	eps := getEndpoints()
	m := http.NewServeMux()
	makeGetVoteResultsHandler(m, eps, []http1.ServerOption{})

	t.Run(testinfo, func(t *testing.T) {
		// Case 1: everything is good, response is http.StatusOK and full struct VodeData;
		u := "/votes/1/results"
		req := httptest.NewRequest(http.MethodGet, u, nil)

		w := httptest.NewRecorder()
		if w == nil || req == nil {
			t.Fatal("w (or r) is nil")
		}

		m.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s (case # 1) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodGet, u, http.StatusOK, resp.StatusCode)
		}

		ep := endpoint.GetVoteResultsResponse{}
		err := json.NewDecoder(resp.Body).Decode(&ep)
		if err != nil {
			t.Errorf("%s, %v", testinfo, err)
		}

		if ep.V0.Header != "Test" {
			t.Errorf("%s (case # 1) failed, expected header \"Test\" but was %s ", testinfo, ep.V0.Header)
		}

		// Case 2: bad vote_id = 2, response must be http.StatusNotFound;
		u = "/votes/2/results"
		req = httptest.NewRequest(http.MethodGet, u, nil)
		w = httptest.NewRecorder()

		m.ServeHTTP(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("%s (case # 2) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodGet, u, http.StatusNotFound, resp.StatusCode)
		}

		// Case 3: bad vote_id = 3, response must be http.StatusInternalServerError;
		u = "/votes/3/results"
		req = httptest.NewRequest(http.MethodGet, u, nil)
		w = httptest.NewRecorder()

		m.ServeHTTP(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("%s (case # 3) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodGet, u, http.StatusInternalServerError, resp.StatusCode)
		}

		// Case 4: bad method POST instead of GET, but vote_id is good, response must be http.StatusMethodNotAllowed;
		u = "/votes/1/results"
		req = httptest.NewRequest(http.MethodPost, u, nil)
		w = httptest.NewRecorder()

		m.ServeHTTP(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("%s (case # 4) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodPost, u, http.StatusMethodNotAllowed, resp.StatusCode)
		}
	})
}

///////////////////////////////////////////
//
// TEST HTTP TRANSPORT UPDATE VOTE RESULTS
//
///////////////////////////////////////////

func TestHttpTransportUpdateVoteResults(t *testing.T) {
	testinfo := "test # 3: UpdateVoteResults"
	eps := getEndpoints()
	m := http.NewServeMux()
	makeUpdateVoteResultsHandler(m, eps, []http1.ServerOption{})

	t.Run(testinfo, func(t *testing.T) {
		// Case 1: everything is good, response must be http.StatusOK;
		u := "/votes"
		dataToSend := &VoteUpdateDTO{
			VoteId:      1,
			ContenderId: 1,
			UserId:      GOOD_USER_ID,
		}
		jsonData, err := json.Marshal(dataToSend)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPut, u, bytes.NewBuffer(jsonData))

		w := httptest.NewRecorder()
		if w == nil {
			t.Fatal("w = nil")
		}

		m.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s (case # 1) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodPut, u, http.StatusOK, resp.StatusCode)
		}

		// Case 2: VoteId is good, but UserId is bad, response must be http.StatusForbidden;
		u = "/votes"
		dataToSend = &VoteUpdateDTO{
			VoteId:      1,
			ContenderId: 1,
			UserId:      BAD_USER_ID,
		}
		jsonData, err = json.Marshal(dataToSend)
		if err != nil {
			t.Fatal(err)
		}

		req = httptest.NewRequest(http.MethodPut, u, bytes.NewBuffer(jsonData))
		w = httptest.NewRecorder()
		m.ServeHTTP(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("%s (case # 2) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodPut, u, http.StatusForbidden, resp.StatusCode)
		}

		// Case 3: VoteId is bad, response must be http.StatusBadRequest;
		u = "/votes"
		dataToSend = &VoteUpdateDTO{
			VoteId:      2,
			ContenderId: 1,
			UserId:      BAD_USER_ID,
		}
		jsonData, err = json.Marshal(dataToSend)
		if err != nil {
			t.Fatal(err)
		}

		req = httptest.NewRequest(http.MethodPut, u, bytes.NewBuffer(jsonData))
		w = httptest.NewRecorder()
		m.ServeHTTP(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("%s (case # 3) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodPut, u, http.StatusBadRequest, resp.StatusCode)
		}

		// Case 4: something is bad, unknown error, response must be http.StatusInternalServerError;
		u = "/votes"
		dataToSend = &VoteUpdateDTO{
			VoteId:      21,
			ContenderId: 1,
			UserId:      BAD_USER_ID,
		}
		jsonData, err = json.Marshal(dataToSend)
		if err != nil {
			t.Fatal(err)
		}

		req = httptest.NewRequest(http.MethodPut, u, bytes.NewBuffer(jsonData))
		w = httptest.NewRecorder()
		m.ServeHTTP(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("%s (case # 4) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodPut, u, http.StatusInternalServerError, resp.StatusCode)
		}

		// Case 5: bad method, POST instead of PUT, response must be http.StatusMethodNotAllowed;
		u = "/votes"
		dataToSend = &VoteUpdateDTO{
			VoteId:      1,
			ContenderId: 1,
			UserId:      GOOD_USER_ID,
		}
		jsonData, err = json.Marshal(dataToSend)
		if err != nil {
			t.Fatal(err)
		}

		req = httptest.NewRequest(http.MethodPost, u, bytes.NewBuffer(jsonData))
		w = httptest.NewRecorder()
		m.ServeHTTP(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("%s (case # 5) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodPut, u, http.StatusMethodNotAllowed, resp.StatusCode)
		}

		// Case 6: bad method, GET instead of PUT, response must be http.StatusMethodNotAllowed;
		u = "/votes"
		req = httptest.NewRequest(http.MethodGet, u, nil)
		w = httptest.NewRecorder()
		m.ServeHTTP(w, req)
		resp = w.Result()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("%s (case # 6) failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodPut, u, http.StatusMethodNotAllowed, resp.StatusCode)
		}
	})
}

//////////////////////////////////////////
//
// TEST HTTP TRANSPORT GET SERVICE STATUS
//
//////////////////////////////////////////

func TestHttpTransportGetServiceStatusOk(t *testing.T) {
	testinfo := "test # 4: GetServiceStatus: Ok"
	ServiceStatus = 0
	eps := getEndpoints()
	m := http.NewServeMux()
	makeGetServiceStatusHandler(m, eps, []http1.ServerOption{})

	t.Run(testinfo, func(t *testing.T) {
		u := "/health"
		req := httptest.NewRequest(http.MethodGet, u, nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodGet, u, http.StatusOK, resp.StatusCode)
		}
	})
}

func TestHttpTransportGetServiceStatusLowMem(t *testing.T) {
	testinfo := "test # 5: GetServiceStatus: Low Memory"
	ServiceStatus = 1
	eps := getEndpoints()
	m := http.NewServeMux()
	makeGetServiceStatusHandler(m, eps, []http1.ServerOption{})

	t.Run(testinfo, func(t *testing.T) {
		u := "/health"
		req := httptest.NewRequest(http.MethodGet, u, nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("%s failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodGet, u, http.StatusServiceUnavailable, resp.StatusCode)
		}
	})
}

func TestHttpTransportGetServiceStatusNoDb(t *testing.T) {
	testinfo := "test # 6: GetServiceStatus: No Database"
	ServiceStatus = 10
	eps := getEndpoints()
	m := http.NewServeMux()
	makeGetServiceStatusHandler(m, eps, []http1.ServerOption{})

	t.Run(testinfo, func(t *testing.T) {
		u := "/health"
		req := httptest.NewRequest(http.MethodGet, u, nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("%s failed, %s %s: expected %d, but was %d",
				testinfo, http.MethodGet, u, http.StatusServiceUnavailable, resp.StatusCode)
		}
	})
}

// --- END ---
