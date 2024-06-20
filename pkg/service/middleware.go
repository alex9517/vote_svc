//  Created : 2024-Mar-14
// Modified : 2024-Mar-27

package service

import (
	"context"

	log "github.com/go-kit/log"
)

// Middleware describes a service middleware.
type Middleware func(VoteService) VoteService

type loggingMiddleware struct {
	logger log.Logger
	next   VoteService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a VoteService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next VoteService) VoteService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) GetVoteData(ctx context.Context, vote_id int) (v0 *VoteData, err error) {
	defer func() {
		l.logger.Log("method", "GetVoteData", "vote_id", vote_id, "v0", v0, "err", err)
	}()
	return l.next.GetVoteData(ctx, vote_id)
}

/*
func (l loggingMiddleware) GetVoteResults(ctx context.Context, vote_id int) (v0 *VoteData, e1 error) {
	defer func() {
		l.logger.Log("method", "GetVoteResults", "vote_id", vote_id, "v0", v0, "e1", e1)
	}()
	return l.next.GetVoteResults(ctx, vote_id)
}
*/

func (l loggingMiddleware) UpdateVoteResults(ctx context.Context, vote_id int, co_id int16, user_id string) (e0 error) {
	defer func() {
		l.logger.Log("method", "UpdateVoteResults", "vote_id", vote_id, "co_id", co_id, "user_id", user_id, "e0", e0)
	}()
	return l.next.UpdateVoteResults(ctx, vote_id, co_id, user_id)
}

func (l loggingMiddleware) GetServiceStatus(ctx context.Context) (v0 *HealthStatus) {
	defer func() {
		l.logger.Log("method", "GetServiceStatus", "v0", v0)
	}()
	return l.next.GetServiceStatus(ctx)
}

// --- END OF FILE ---
