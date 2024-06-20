//  Created : 2024-Mar-27
// Modified : 2024-Apr-24

package endpoint

import (
	"context"
	"net/http"
	"testing"
	"time"
	"vote_svc/pkg/service"

	"github.com/patrickmn/go-cache"
)

const TESTDATA_HEADER = "Some Competition 2024"
const TESTDATA_MESSAGE = "Vote for the best candidate"
const TESTDATA_DEADLINE = "2030-04-29T22:00:00Z"
const TESTDATA_USER_ID = "5287eb7d-4813-49e7-8e2d-c8a6ce8b3c4c"
const TESTDATA_CO_INFO = "Good Person "

var voteGetVoteDataMock func(ctx context.Context, vote_id int) (*service.VoteData, error)
var voteUpdateVoteResultsMock func(ctx context.Context, vote_id int, co_id int16, user_id string) error
var voteGetServiceStatusMock func(ctx context.Context) *service.HealthStatus

type voteServiceMock struct{}

func (b voteServiceMock) GetVoteData(ctx context.Context, vote_id int) (*service.VoteData, error) {
	return voteGetVoteDataMock(ctx, vote_id)
}

func (b voteServiceMock) UpdateVoteResults(ctx context.Context, vote_id int, co_id int16, user_id string) error {
	return voteUpdateVoteResultsMock(ctx, vote_id, co_id, user_id)
}

func (b voteServiceMock) GetServiceStatus(ctx context.Context) *service.HealthStatus {
	return voteGetServiceStatusMock(ctx)
}

func newServiceMock([]service.Middleware) service.VoteService {
	return &voteServiceMock{}
}

////////////////////////////////////
//
// TEST MAKE GET VOTE DATA ENDPOINT
//
////////////////////////////////////

func TestMakeGetVoteDataEndpoint(t *testing.T) {
	testinfo := "test # 1: GetVoteData endpoint"
	memCache := cache.New(5*time.Second, 10*time.Second)
	svc := newServiceMock([]service.Middleware{})
	endpoint := MakeGetVoteDataEndpoint(svc, memCache)
	var deadline, _ = time.Parse(time.RFC3339, TESTDATA_DEADLINE)

	t.Run(testinfo, func(t *testing.T) {
		voteGetVoteDataMock = func(_ context.Context, vote_id int) (*service.VoteData, error) {
			if vote_id == 1 {
				contenders := make([]service.Contender, 2)
				contenders[0] = service.Contender{
					Id:      1,
					Name:    "Alex Good",
					Alias:   "Alex # 1",
					Info:    TESTDATA_CO_INFO,
					Count:   0,
					Updated: time.Now(),
					Picture: "AlexGood.jpg",
				}
				contenders[1] = service.Contender{
					Id:      2,
					Name:    "Alex Bravo",
					Alias:   "Alex # 2",
					Info:    TESTDATA_CO_INFO,
					Count:   0,
					Updated: time.Now(),
					Picture: "AlexBravo.jpg",
				}

				return &service.VoteData{
					VoteId:       1,
					Header:       TESTDATA_HEADER,
					Message:      TESTDATA_MESSAGE,
					Deadline:     deadline,
					Authenticate: false,
					AllowResults: true,
					Contenders:   contenders,
				}, nil
			} else {
				return nil, service.ErrNotFound
			}
		}

		// Case 1: vote_id = 1, service returns good 'VoteData';
		vote_id := 1
		r, _ := endpoint(context.Background(), GetVoteDataRequest{VoteId: vote_id})
		if v, ok := r.(GetVoteDataResponse); ok {
			if v.E1 != nil {
				t.Errorf("%v (case # 1) failed, err %v", testinfo, v.E1)
			}
		}

		// Case 2: vote_id = 21, service returns 'ErrNotFound';
		vote_id = 21
		r, _ = endpoint(context.Background(), GetVoteDataRequest{VoteId: vote_id})
		if v, ok := r.(GetVoteDataResponse); ok {
			if v.E1 != service.ErrNotFound {
				t.Errorf("%v (case # 2) failed, err %v, must be %v", testinfo, v.E1, service.ErrNotFound)
			}
		}
	})
}

///////////////////////////////////////
//
// TEST MAKE GET VOTE RESULTS ENDPOINT
//
///////////////////////////////////////

func TestMakeGetVoteResultsEndpoint(t *testing.T) {
	testinfo := "test # 2: GetVoteResults endpoint"
	memCache := cache.New(5*time.Second, 10*time.Second)
	svc := newServiceMock([]service.Middleware{})
	endpoint := MakeGetVoteResultsEndpoint(svc, memCache)
	var deadline, _ = time.Parse(time.RFC3339, TESTDATA_DEADLINE)

	t.Run(testinfo, func(t *testing.T) {
		voteGetVoteDataMock = func(_ context.Context, vote_id int) (*service.VoteData, error) {
			if vote_id == 1 {
				contenders := make([]service.Contender, 2)
				contenders[0] = service.Contender{
					Id:      1,
					Name:    "Alex Good",
					Alias:   "Alex # 1",
					Info:    TESTDATA_CO_INFO,
					Count:   101,
					Updated: time.Now(),
					Picture: "AlexGood.jpg",
				}
				contenders[1] = service.Contender{
					Id:      2,
					Name:    "Alex Bravo",
					Alias:   "Alex # 2",
					Info:    TESTDATA_CO_INFO,
					Count:   106,
					Updated: time.Now(),
					Picture: "AlexBravo.jpg",
				}

				return &service.VoteData{
					VoteId:       1,
					Header:       TESTDATA_HEADER,
					Message:      TESTDATA_MESSAGE,
					Deadline:     deadline,
					Authenticate: false,
					AllowResults: true,
					Contenders:   contenders,
				}, nil
			} else {
				return nil, service.ErrNotFound
			}
		}

		// Case 1: vote_id = 1, service returns good 'VoteData';
		vote_id := 1
		r, _ := endpoint(context.Background(), GetVoteResultsRequest{VoteId: vote_id})
		if v, ok := r.(GetVoteResultsResponse); ok {
			if v.E1 != nil {
				t.Errorf("%v (case # 1) failed, err %v", testinfo, v.E1)
			}
		}

		// Case 2: vote_id = 21, service returns 'ErrNotFound';
		vote_id = 21
		r, _ = endpoint(context.Background(), GetVoteResultsRequest{VoteId: vote_id})
		if v, ok := r.(GetVoteResultsResponse); ok {
			if v.E1 != service.ErrNotFound {
				t.Errorf("%v (case # 2) failed, err %v, must be %v", testinfo, v.E1, service.ErrNotFound)
			}
		}
	})
}

//////////////////////////////////////////
//
// TEST MAKE UPDATE VOTE RESULTS ENDPOINT
//
//////////////////////////////////////////

func TestMakeUpdateVoteResultsEndpoint(t *testing.T) {
	testinfo := "test # 3: UpdateVoteResults endpoint"
	svc := newServiceMock([]service.Middleware{})
	endpoint := MakeUpdateVoteResultsEndpoint(svc)
	// var deadline, _ = time.Parse(time.RFC3339, TESTDATA_DEADLINE)

	t.Run(testinfo, func(t *testing.T) {
		voteUpdateVoteResultsMock = func(_ context.Context, vote_id int, co_id int16, user_id string) error {
			if user_id == TESTDATA_USER_ID {
				return service.ErrForbidden
			} else if vote_id == 1 && co_id == 1 {
				return nil
			} else if vote_id == 2 {
				return service.ErrBadRequest
			} else {
				return service.ErrInternalServerError
			}
		}
		// Case 1: vote_id = 1, co_id = 1, user_id != TESTDATA_USER_ID, service returns 'nil';
		good_user_id := "a8597900-9aa0-40d9-9dcc-ff1f4210d7d8"
		r, _ := endpoint(context.Background(), UpdateVoteResultsRequest{VoteId: 1, ContenderId: 1, UserId: good_user_id})
		if v, ok := r.(UpdateVoteResultsResponse); ok {
			if v.E0 != nil {
				t.Errorf("%v (case # 1) failed, err %v (must be nil)", testinfo, v.E0)
			}
		}

		// Case 2: vote_id = 1, co_id = 1, user_id == TESTDATA_USER_ID, service returns ErrForbidden;
		r, _ = endpoint(context.Background(), UpdateVoteResultsRequest{VoteId: 1, ContenderId: 1, UserId: TESTDATA_USER_ID})
		if v, ok := r.(UpdateVoteResultsResponse); ok {
			if v.E0 != service.ErrForbidden {
				t.Errorf("%v (case # 2) failed, err %v (must be %v)", testinfo, v.E0, service.ErrForbidden)
			}
		}

		// Case 3: vote_id = 2, co_id = 1, user_id != TESTDATA_USER_ID, service returns ErrBadRequest;
		r, _ = endpoint(context.Background(), UpdateVoteResultsRequest{VoteId: 2, ContenderId: 1, UserId: good_user_id})
		if v, ok := r.(UpdateVoteResultsResponse); ok {
			if v.E0 != service.ErrBadRequest {
				t.Errorf("%v (case # 3) failed, err %v (must be %v)", testinfo, v.E0, service.ErrBadRequest)
			}
		}

		// Case 4: vote_id = 3, co_id = 1, user_id != TESTDATA_USER_ID, service returns ErrInternalServerError;
		r, _ = endpoint(context.Background(), UpdateVoteResultsRequest{VoteId: 3, ContenderId: 1, UserId: good_user_id})
		if v, ok := r.(UpdateVoteResultsResponse); ok {
			if v.E0 != service.ErrInternalServerError {
				t.Errorf("%v (case # 4) failed, err %v (must be %v)", testinfo, v.E0, service.ErrInternalServerError)
			}
		}
	})
}

/////////////////////////////////////////
//
// TEST MAKE GET SERVICE STATUS ENDPOINT
//
/////////////////////////////////////////

func TestMakeGetServiceStatusEndpoint(t *testing.T) {
	testinfo := "test # 4: GetServiceStatus endpoint"
	svc := newServiceMock([]service.Middleware{})
	endpoint := MakeGetServiceStatusEndpoint(svc)

	t.Run(testinfo, func(t *testing.T) {
		// Case # 1: service returns good 'HealthStatus';
		voteGetServiceStatusMock = func(_ context.Context) *service.HealthStatus {
			return &service.HealthStatus{
				Status:  http.StatusOK,
				Message: service.SERVICE_STATUS_OK,
			}
		}

		r, _ := endpoint(context.Background(), GetServiceStatusRequest{})
		if v, ok := r.(GetServiceStatusResponse); ok {
			if v.H0.Status != http.StatusOK {
				t.Errorf("%v (case # 1) failed, status %v, must be %v", testinfo, v.H0.Status, http.StatusOK)
			}

			if v.H0.Message != service.SERVICE_STATUS_OK {
				t.Errorf("%v (case # 1) failed, message %v, must be %v", testinfo, v.H0.Message, service.SERVICE_STATUS_OK)
			}
		}
		// Case # 2: service returns 'ServiceUnavailable, low memory';
		voteGetServiceStatusMock = func(_ context.Context) *service.HealthStatus {
			return &service.HealthStatus{
				Status:  http.StatusServiceUnavailable,
				Message: service.SERVICE_STATUS_LOW_MEMORY,
			}
		}

		r, _ = endpoint(context.Background(), GetServiceStatusRequest{})
		if v, ok := r.(GetServiceStatusResponse); ok {
			if v.H0.Status != http.StatusServiceUnavailable {
				t.Errorf("%v (case # 2) failed, status %v, must be %v", testinfo, v.H0.Status, http.StatusServiceUnavailable)
			}

			if v.H0.Message != service.SERVICE_STATUS_LOW_MEMORY {
				t.Errorf("%v (case # 2) failed, message %v, must be %v", testinfo, v.H0.Message, service.SERVICE_STATUS_LOW_MEMORY)
			}
		}
		// Case # 3: service returns 'ServiceUnavailable, no database';
		voteGetServiceStatusMock = func(_ context.Context) *service.HealthStatus {
			return &service.HealthStatus{
				Status:  http.StatusServiceUnavailable,
				Message: service.SERVICE_STATUS_NO_DATABASE,
			}
		}

		r, _ = endpoint(context.Background(), GetServiceStatusRequest{})
		if v, ok := r.(GetServiceStatusResponse); ok {
			if v.H0.Status != http.StatusServiceUnavailable {
				t.Errorf("%v (case # 3) failed, status %v, must be %v", testinfo, v.H0.Status, http.StatusServiceUnavailable)
			}

			if v.H0.Message != service.SERVICE_STATUS_NO_DATABASE {
				t.Errorf("%v (case # 3) failed, message %v, must be %v", testinfo, v.H0.Message, service.SERVICE_STATUS_NO_DATABASE)
			}
		}
	})
}

// --- END OF FILE ---
