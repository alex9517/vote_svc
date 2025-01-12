// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
// ↑
// Don't bother about it. In most cases this file needs editing.

//  Created : 2024-Mar-14
// Modified : 2024-Apr-04

package http

import (
	kithttp "github.com/go-kit/kit/transport/http"
	http "net/http"
	endpoint "vote_svc/pkg/endpoint"
	"os"
	"strings"
	"github.com/rs/cors"
)

// NewHTTPHandler returns a handler that makes a set of endpoints available on predefined paths.

func NewHTTPHandler(endpoints endpoint.Endpoints, options map[string][]kithttp.ServerOption) http.Handler {
	m := http.NewServeMux()
	makeGetVoteDataHandler(m, endpoints, options["GetVoteData"])
	makeGetVoteResultsHandler(m, endpoints, options["GetVoteResults"])
	makeUpdateVoteResultsHandler(m, endpoints, options["UpdateVoteResults"])
	makeGetServiceStatusHandler(m, endpoints, options["GetServiceStatus"])

	// CORS-related stuff (Cross-Origin Resource Sharing).
	// This was not auto generated, but it's required;

	s := os.Getenv("ALLOW_ORIGINS")
	if s == "" {
		s = "http://localhost"
	}

	AllowOrigins := strings.Split(s, ",")

	cors := cors.New(cors.Options{
		AllowedOrigins: AllowOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			// http.MethodOptions,
			// http.MethodHead,
		},
		MaxAge: 15,
		AllowedHeaders: []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
		OptionsPassthrough: false,
		Debug: true,
	})
	
	handler := cors.Handler(m)
	return handler
}

// --- END ---
