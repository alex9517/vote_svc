//  Created : 2024-Apr-04
// Modified : 2024-Jun-18

// The following marks in the func header comment mean:
//    +++        func can be used in other projects without change.
//    ++-        func can be used in other projects as is, or may need some editing.
//    ---        func requires modifications to comply to proj specifics.

package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	endpoint "vote_svc/pkg/endpoint"
	pkghttp "vote_svc/pkg/http"
	service "vote_svc/pkg/service"

	kitendpoint "github.com/go-kit/kit/endpoint"
	prometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/ratelimit"
	opentracing "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	log "github.com/go-kit/log"
	"github.com/gocql/gocql"
	consulapi "github.com/hashicorp/consul/api"
	lightsteptracergo "github.com/lightstep/lightstep-tracer-go"
	group "github.com/oklog/oklog/pkg/group"
	opentracinggo "github.com/opentracing/opentracing-go"
	zipkingoopentracing "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkingo "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/patrickmn/go-cache"
	prometheus1 "github.com/prometheus/client_golang/prometheus"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	appdash "sourcegraph.com/sourcegraph/appdash"
	opentracing1 "sourcegraph.com/sourcegraph/appdash/opentracing"
)

const DEFAULT_DEBUG_ADDR = ":8080"
const DEFAULT_HTTP_ADDR = ":8081"
const DEFAULT_HTTPS_ADDR = ":8443"

// The default rate limit for 'GET' requests is 300 req/sec.
// It can be changed by the 'rate-limit' cmdline param.
const DEFAULT_RATE_LIMIT = 300

// The short burst [of requests] limit = ratelimit * factor.
const RATE_BURST_FACTOR = 4

// The default rate limit for 'PUT' requests.
const UPDATE_RATE_LIMIT = DEFAULT_RATE_LIMIT / 5

// The rate limit for '/health' endpoint.
// Why would anybody send many request to this endpoint?
const HEALTH_RATE_LIMIT = 2

// This is Apache Cassandra Database entry node.
// It can be overwritten by the cmdline param 'database-url'.
const DEFAULT_DATABASE_URL = "172.16.70.31"

// Keyspace name must coincide with the actual keyspace in the database.
const DEFAULT_DATABASE_KEYSPACE = "polls"

// Most likely you won't need to change this port.
const DEFAULT_DATABASE_PORT = 9042

// Cache.
const DEFAULT_CACHE_EXPIRE = 15 * time.Minute
const DEFAULT_CACHE_CLEAR = 30 * time.Minute

// Certificates [ for Docker image ]
const SERVICE_CERT = "/etc/x509/https/service.pem" // See also cmdline flags;
const SERVICE_KEY = "/etc/x509/https/service.key"  // See ...

// Consul (service discovery/registration), if that Consul is avaialble.
const REG_SERVICE_ID = "vote-svc"
const REG_SERVICE_NAME = "Votes microservice"
const REG_SERVICE_INTERVAL = "60s"
const REG_SERVICE_TIMEOUT = "30s"

// It's external ...
var tracer opentracinggo.Tracer

// The one and the only.
var logger log.Logger

// For the cmdline flags 'vote-svc' is just a FlagSet name; probably it does not
// affect the app's behavior, but it will be displayed in the usage and error messages;
var fs = flag.NewFlagSet("vote-svc", flag.ExitOnError)

var debugAddr = fs.String("debug-addr", DEFAULT_DEBUG_ADDR, "Debug and metrics listen address")
var httpAddr = fs.String("http-addr", DEFAULT_HTTP_ADDR, "HTTP listen address")
var httpsAddr = fs.String("https-addr", DEFAULT_HTTPS_ADDR, "HTTPS listen address")
var serviceCert = fs.String("service-cert", SERVICE_CERT, "Certificate file for TLS")
var serviceKey = fs.String("service-key", SERVICE_KEY, "Private key file for TLS")
var zipkinURL = fs.String("zipkin-url", "", "Enable Zipkin tracing via a collector URL e.g. http://localhost:9411/api/v1/spans")
var lightstepToken = fs.String("lightstep-token", "", "Enable LightStep tracing via a LightStep access token")
var appdashAddr = fs.String("appdash-addr", "", "Enable Appdash tracing via an Appdash server host:port")
var rateLimit = fs.Int("rate-limit", DEFAULT_RATE_LIMIT, "Rate limit for requests")
var rateLimitPut = fs.Int("rate-limit-put", UPDATE_RATE_LIMIT, "Rate limit for PUT/POST requests")
var databaseURL = fs.String("database-url", DEFAULT_DATABASE_URL, "Database URL (whatever it means for the database you use)")
var cacheExpire = fs.Duration("cache-expire", DEFAULT_CACHE_EXPIRE, "Cache expiration time")
var cacheClear = fs.Duration("cache-clear", DEFAULT_CACHE_CLEAR, "Cache clear time")

// var databasePort = fs.Int("database-port", DEFAULT_DATABASE_PORT, "Cassandra Database port")
// var databaseKeyspace = fs.String("database-keyspace", DEFAULT_DATABASE_KEYSPACE, "Cassandra Database keyspace")

///////////
//
// M A I N
//
/////////// ++-

func main() {

	fs.Parse(os.Args[1:])

	// Create a single logger that will be used here and given to other components.
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// Determine which tracer to use. It will be passed as a dependency
	// to all the components using it.
	if *zipkinURL != "" {
		logger.Log("tracer", "Zipkin", "URL", *zipkinURL)
		reporter := zipkinhttp.NewReporter(*zipkinURL)
		defer reporter.Close()
		endpoint, err := zipkingo.NewEndpoint("vote-svc", "localhost:80") // vote-svc ?
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		localEndpoint := zipkingo.WithLocalEndpoint(endpoint)
		nativeTracer, err := zipkingo.NewTracer(reporter, localEndpoint)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		tracer = zipkingoopentracing.Wrap(nativeTracer)
	} else if *lightstepToken != "" {
		logger.Log("tracer", "LightStep")
		tracer = lightsteptracergo.NewTracer(lightsteptracergo.Options{AccessToken: *lightstepToken})
		defer lightsteptracergo.Flush(context.Background(), tracer)
	} else if *appdashAddr != "" {
		logger.Log("tracer", "Appdash", "addr", *appdashAddr)
		collector := appdash.NewRemoteCollector(*appdashAddr)
		tracer = opentracing1.NewTracer(collector)
		defer collector.Close()
	} else {
		logger.Log("tracer", "none")
		tracer = opentracinggo.GlobalTracer()
	}

	// Database (Apache Cassandra noSQL). Technically, 'cluster' is just a pointer
	// to a struct 'gocql.ClusterConfig'. The actual connection to the database is
	// established in the 'service' layer (see 'pkg/service/service.go'). Notice that
	// this app does not use password auth for the database access, but it can be changed.
	cluster := gocql.NewCluster(*databaseURL)
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 4
	cluster.Timeout = time.Second * 10
	cluster.ConnectTimeout = time.Second * 10
	cluster.ReconnectionPolicy = &gocql.ConstantReconnectionPolicy{MaxRetries: 10, Interval: 6 * time.Second}

	// cluster.ProtoVersion = 4
	// cluster.Keyspace = databaseKeyspace
	// cluster.Port = databasePort
	// cluster.Hosts = []string{*databaseURL}
	// cluster.PoolConfig.HostSelectionPolicy = gocql.HostPoolHostPolicy(hostpool.New(nil))

	// cluster.Authenticator = gocql.PasswordAuthenticator{
	// Username: "user",
	// Password: "password",
	// }

	// Note! This is not memcached! This is local in-memory cache.
	memCache := cache.New(*cacheExpire, *cacheClear)

	// Register [our] service with 'Consul' (if 'Consul' service is available).
	registerServiceWithConsul()

	//
	// ===== Main part =====
	//
	svc := service.New(cluster, getServiceMiddleware(logger))
	eps := endpoint.New(svc, memCache, getEndpointMiddleware(logger))
	g := createService(eps)
	initMetricsEndpoint(g)
	initCancelInterrupt(g)
	logger.Log("exit", g.Run())
}

//////////////////
//
// CREATE SERVICE
//
////////// called by main ++-

func createService(endpoints endpoint.Endpoints) (g *group.Group) {
	g = &group.Group{}
	initHttpHandler(endpoints, g)
	return g
}

/////////////////////
//
// INIT HTTP HANDLER
//
////////// called by createService ++-

func initHttpHandler(endpoints endpoint.Endpoints, g *group.Group) {
	options := defaultHttpOptions(logger, tracer)

	// Add your http options here.

	httpHandler := pkghttp.NewHTTPHandler(endpoints, options)
	httpListener, err := getHttpListenerWithTLS(*httpsAddr)
	if err != nil {
		// Failed to get listener with TLS,
		// let's make it old style (trivial HTTP without TLS).
		httpListener, err = net.Listen("tcp", *httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "addr", *httpAddr, "during", "Listen", "err", err)
		} else {
			logger.Log("transport", "HTTP", "addr", *httpAddr)
		}
	} else {
		// We're gonna use HTTPS.
		logger.Log("transport", "HTTPS", "addr", *httpsAddr)
	}

	g.Add(func() error {
		return http.Serve(httpListener, httpHandler)
	}, func(error) {
		httpListener.Close()
	})
}

//////////////////////////////
//
// GET HTTP LISTENER WITH TLS
//
////////// called by initHttpHandler +++

func getHttpListenerWithTLS(addr string) (net.Listener, error) {
	if addr == "" {
		err := errors.New("TCP port for HTTPS is not specified")
		logger.Log("transport", "HTTPS",
			"addr", addr, "during", "getHttpListenerWithTLS", "err", err)
		return nil, err
	}

	//  Load our [ self-signed ] certificate.
	cert, err := tls.LoadX509KeyPair(*serviceCert, *serviceKey)
	if err != nil {
		logger.Log("transport", "HTTPS", "addr", addr, "during", "LoadX509KeyPair", "err", err)
		return nil, err
	}

	// Certificate seems good, let's try to create a listener.
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	httpListener, err := tls.Listen("tcp", addr, config)
	if err != nil {
		logger.Log("transport", "HTTPS", "addr", addr, "during", "Listen", "err", err)
		return nil, err
	}

	return httpListener, nil
}

////////////////////////
//
// DEFAULT HTTP OPTIONS
//
////////// called by initHttpHandler ---

func defaultHttpOptions(logger log.Logger, tracer opentracinggo.Tracer) map[string][]kithttp.ServerOption {
	options := map[string][]kithttp.ServerOption{
		"GetVoteData":       {kithttp.ServerErrorEncoder(pkghttp.ErrorEncoder), kithttp.ServerErrorLogger(logger), kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "GetVoteData", logger))},
		"GetVoteResults":    {kithttp.ServerErrorEncoder(pkghttp.ErrorEncoder), kithttp.ServerErrorLogger(logger), kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "GetVoteResults", logger))},
		"UpdateVoteResults": {kithttp.ServerErrorEncoder(pkghttp.ErrorEncoder), kithttp.ServerErrorLogger(logger), kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "UpdateVoteResults", logger))},
		"GetServiceStatus":  {kithttp.ServerErrorEncoder(pkghttp.ErrorEncoder), kithttp.ServerErrorLogger(logger), kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "GetServiceStatus", logger))},
	}
	return options
}

/////////////////////////
//
// GET SERVICE MIDDLEWARE
//
////////// called by main ++-

func getServiceMiddleware(logger log.Logger) (mw []service.Middleware) {
	mw = []service.Middleware{}
	mw = addDefaultServiceMiddleware(logger, mw)

	// Append your middleware here

	return
}

//////////////////////////////////
//
// ADD DEFAULT SERVICE MIDDLEWARE
//
////////// called by getServiceMiddleware +++

func addDefaultServiceMiddleware(logger log.Logger, mw []service.Middleware) []service.Middleware {
	return append(mw, service.LoggingMiddleware(logger))
}

///////////////////////////
//
// GET ENDPOINT MIDDLEWARE
//
//////////// called by main ++-

func getEndpointMiddleware(logger log.Logger) (mw map[string][]kitendpoint.Middleware) {
	mw = map[string][]kitendpoint.Middleware{}
	duration := prometheus.NewSummaryFrom(prometheus1.SummaryOpts{
		Help:      "Request duration in seconds.",
		Name:      "request_duration_seconds",
		Namespace: "example",
		Subsystem: "vote_svc", // NOTE! Setting {Subsystem: "vote-svc"}, it breaks app;
	}, []string{"method", "success"})

	addDefaultEndpointMiddleware(logger, duration, mw)

	// Add you endpoint middleware here

	return
}

///////////////////////////////////
//
// ADD DEFAULT ENDPOINT MIDDLEWARE
//
////////// called by getEndpointMiddleware ---

func addDefaultEndpointMiddleware(logger log.Logger, duration *prometheus.Summary, mw map[string][]kitendpoint.Middleware) {
	mw["GetVoteData"] = []kitendpoint.Middleware{
		endpoint.LoggingMiddleware(log.With(logger, "method", "GetVoteData")),
		endpoint.InstrumentingMiddleware(duration.With("method", "GetVoteData")),
		ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(*rateLimit), RATE_BURST_FACTOR*(*rateLimit)))}

	mw["GetVoteResults"] = []kitendpoint.Middleware{
		endpoint.LoggingMiddleware(log.With(logger, "method", "GetVoteResults")),
		endpoint.InstrumentingMiddleware(duration.With("method", "GetVoteResults")),
		ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(*rateLimit), RATE_BURST_FACTOR*(*rateLimit)))}

	mw["UpdateVoteResults"] = []kitendpoint.Middleware{
		endpoint.LoggingMiddleware(log.With(logger, "method", "UpdateVoteResults")),
		endpoint.InstrumentingMiddleware(duration.With("method", "UpdateVoteResults")),
		ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(*rateLimitPut), RATE_BURST_FACTOR*(*rateLimitPut)))}

	mw["GetServiceStatus"] = []kitendpoint.Middleware{
		endpoint.LoggingMiddleware(log.With(logger, "method", "GetServiceStatus")),
		endpoint.InstrumentingMiddleware(duration.With("method", "GetServiceStatus")),
		ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(HEALTH_RATE_LIMIT), RATE_BURST_FACTOR*(HEALTH_RATE_LIMIT)))}
}

/////////////////////////
//
// INIT METRICS ENDPOINT
//
////////// called by main +++

func initMetricsEndpoint(g *group.Group) {
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
	debugListener, err := net.Listen("tcp", *debugAddr)
	if err != nil {
		logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "debug/HTTP", "addr", *debugAddr)
		return http.Serve(debugListener, http.DefaultServeMux)
	}, func(error) {
		debugListener.Close()
	})
}

/////////////////////////
//
// INIT CANCEL INTERRUPT
//
////////// called by main +++

func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}

////////////////////////////////
//
// REGISTER SERVICE WITH CONSUL
//
///////////////// called by main +++

func registerServiceWithConsul() {
	conf := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(conf)
	if err != nil {
		logger.Log("main", "Service Registration", "err", err)
	}

	port_num, err := strconv.Atoi(strings.TrimPrefix(*httpAddr, ":"))
	if err != nil {
		logger.Log("main", "Service Registartion, port num", "err", err)
	}

	port := port_num
	address, err := os.Hostname()
	if err != nil {
		logger.Log("main", "Service Registartion, hostname", "err", err)
	}

	registration := &consulapi.AgentServiceRegistration{
		ID:      REG_SERVICE_ID,
		Name:    REG_SERVICE_NAME,
		Port:    port,
		Address: address,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s%s/check", address, *httpAddr),
			Interval: REG_SERVICE_INTERVAL,
			Timeout:  REG_SERVICE_TIMEOUT,
		},
	}

	regiErr := consul.Agent().ServiceRegister(registration)

	if regiErr != nil {
		logger.Log("main", "Service Registartion failure", "err", regiErr)
	} else {
		logger.Log("main", "Service Registartion success", "gRPC", address+*httpAddr)
	}
}

// --- END OF FILE ---
