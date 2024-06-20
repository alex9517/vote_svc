//  Created : 2024-Mar-14
// Modified : 2024-May-14

package service

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

// These are HealthStatus messages.
const SERVICE_STATUS_OK = "UP"
const SERVICE_STATUS_DOWN = "DOWN"
const SERVICE_STATUS_UNKNOWN = "Unknown"
const SERVICE_STATUS_LOW_MEMORY = "Unsufficient memory"
const SERVICE_STATUS_NO_DATABASE = "Database connection failure"

// The following error messages should not be capitalized (compiler warning);
const ERR_MSG_NO_CONTENT = "no content"
const ERR_MSG_BAD_REQUEST = "bad request"
const ERR_MSG_UNAUTHORIZED = "unauthorized"
const ERR_MSG_FORBIDDEN = "forbidden"
const ERR_MSG_NOT_FOUND = "not found"
const ERR_MSG_METHOD_NOT_ALLOWED = "method not allowed"
const ERR_MSG_UNAVAILABLE = "service unavailable"
const ERR_MSG_SERVER_ERROR = "internal server error"

const LOW_MEM_THRESHOLD uint64 = 1048576 // Mem size in KB (~1GB);

var (
	ErrNoContent           = errors.New(ERR_MSG_NO_CONTENT)
	ErrBadRequest          = errors.New(ERR_MSG_BAD_REQUEST)
	ErrUnauthorized        = errors.New(ERR_MSG_UNAUTHORIZED)
	ErrForbidden           = errors.New(ERR_MSG_FORBIDDEN)
	ErrNotFound            = errors.New(ERR_MSG_NOT_FOUND)
	ErrMethodNotAllowed    = errors.New(ERR_MSG_METHOD_NOT_ALLOWED)
	ErrServiceUnavailable  = errors.New(ERR_MSG_UNAVAILABLE)
	ErrInternalServerError = errors.New(ERR_MSG_SERVER_ERROR)
)

// VoteSvcService describes the service.
type VoteService interface {
	GetVoteData(ctx context.Context, vote_id int) (*VoteData, error)
	UpdateVoteResults(ctx context.Context, vote_id int, co_id int16, user_id string) error
	GetServiceStatus(ctx context.Context) *HealthStatus
}

type Contender struct {
	Id      int16     `json:"id"`
	Name    string    `json:"name"`
	Alias   string    `json:"alias"`   // A shorter version of name;
	Info    string    `json:"info"`    // Compact description, 3..10 lines;
	Picture string    `json:"picture"` // A filename to be added to URL;
	Count   int64     `json:"count"`   // Number of votes for this contender;
	Updated time.Time `json:"updated"` // Last count update timestamp;
}

type VoteData struct {
	VoteId       int         `json:"vote_id"`
	Header       string      `json:"header"`
	Message      string      `json:"message"`
	Resources    string      `json:"resources"` // URL to external resources, e.g. images;
	Deadline     time.Time   `json:"deadline"`
	Authenticate bool        `json:"authenticate"`
	AllowResults bool        `json:"allow_results"`
	Contenders   []Contender `json:"contenders"`
}

type VoteDataAux struct {
	voteId       int
	co_id        int16
	header       string
	message      string
	resources    string
	deadline     time.Time
	authenticate bool
	allowresults bool
	co_name      string
	co_alias     string
	co_info      string
	co_picture   string
	co_count     int64
	co_updated   time.Time
}

type HealthStatus struct {
	Status  int32  `json:"health_status"`
	Message string `json:"health_message"`
}

// The database is Apache Cassandra noSQL/CQL.
type basicVoteService struct {
	db *gocql.ClusterConfig
}

/////////////////
//
// GET VOTE DATA
//
/////////////////

// 1. It provides data required to fill the [ browser client ] form for a specified 'vote_id'.
// 2. It provides results related to a specified 'vote_id', including data required for diagram.

// In both cases this func returns the same struct containing all the data available, and the
// client app decides what to use and what to ignore, e.g. in frist case it ignores counts of
// votes for contenders, while in second case those counts are used to create a graphic diagram.

// The 'VoteData' struct includes info about vote purpose, candidates/contenders, etc. All data is
// fetched from the database. The number of records related to a specific 'vote_id' is equal to the
// number of contenders which is not supposed to be large, maybe 2..20). This data is packed into
// a single struct 'VoteData' containing an array (i.e. slice) of the type struct 'Contenders'
// with the data about contenders (id, name, picture, ...).

func (b *basicVoteService) GetVoteData(ctx context.Context, vote_id int) (*VoteData, error) {

	if b.db == nil {
		return nil, ErrServiceUnavailable
	}

	session, err := b.db.CreateSession()
	if err != nil {
		return nil, err
	}

	defer session.Close()

	records := []VoteDataAux{} // This is an intermediate slice to receive the database records;
	m := map[string]interface{}{}

	// Note that 'stmt' here is a CQL statement (Cassandra Query Language), not SQL.
	// In general, it's supposed to fetch at least two records (two contenders are minimum,
	// there is no sense to have election if you have only one candidate). Once again, this
	// is not SQL database, the tables are not normalized and some data is duplicated.

	stmt := `SELECT vote_id, co_id, header, message, resources, deadline, authenticate,
	 allowresults, co_name, co_alias, co_info, co_picture, co_count, co_updated
	 FROM polls.votes WHERE vote_id = ?`

	// iterable := session.Query(stmt, vote_id).WithContext(ctx).Consistency(gocql.One).Iter()
	iterable := session.Query(stmt, vote_id).Iter()

	for iterable.MapScan(m) {
		records = append(records, VoteDataAux{
			voteId:       m["vote_id"].(int),
			co_id:        m["co_id"].(int16),
			header:       m["header"].(string),
			message:      m["message"].(string),
			resources:    m["resources"].(string),
			deadline:     m["deadline"].(time.Time),
			authenticate: m["authenticate"].(bool),
			allowresults: m["allowresults"].(bool),
			co_name:      m["co_name"].(string),
			co_alias:     m["co_alias"].(string),
			co_info:      m["co_info"].(string),
			co_picture:   m["co_picture"].(string),
			co_count:     m["co_count"].(int64),
			co_updated:   m["co_updated"].(time.Time),
		})

		m = map[string]interface{}{} // Do not remove this!
	}

	if len(records) == 0 {
		return nil, ErrNotFound
	}

	// Let's repack database records into a single struct 'VoteData'.
	var contenders []Contender
	for _, r := range records {
		contenders = append(contenders, Contender{
			Id:      r.co_id,
			Name:    r.co_name,
			Alias:   r.co_alias,
			Info:    r.co_info,
			Count:   r.co_count,
			Updated: r.co_updated,
			Picture: r.co_picture,
		})
	}

	var res VoteData = VoteData{
		VoteId:       records[0].voteId,
		Header:       records[0].header,
		Message:      records[0].message,
		Resources:    records[0].resources,
		Deadline:     records[0].deadline,
		Authenticate: records[0].authenticate,
		AllowResults: records[0].allowresults,
		Contenders:   contenders,
	}

	return &res, nil
}

///////////////////////
//
// UPDATE VOTE RESULTS
//
///////////////////////

// This func performs following ops with the database:

// 1. It checks if the specified 'vote_id' and 'co_id' are valid (i.e. present in the database)
// and the current datetime is before the deadline. If not, it returns an error.

// 2. It tries to insert a new record into the 'voters' table (new 'user_id')
// to prevent this user/voter from voting again. In case of failure, it returns an error.

// 3. It updates the 'votes' table incrementing the 'co_count' of the specified contender.

func (b *basicVoteService) UpdateVoteResults(ctx context.Context, vote_id int, co_id int16, user_id string) (e0 error) {
	if b.db == nil {
		return ErrServiceUnavailable
	}

	session, err := b.db.CreateSession()
	if err != nil {
		return err
	}

	defer session.Close()

	// Step # 1: let's check if the 'vote_id' and 'co_id' are valid and deadline is in the future;
	stmt := "SELECT deadline FROM polls.votes WHERE vote_id = ? AND co_id = ?"

	var deadline time.Time // This would be the number of records found in the 'votes' table;
	err = session.Query(stmt, vote_id, co_id).WithContext(ctx).Scan(&deadline)
	if err != nil { // This is supposed to be ErrNotFound if CQL returns 0 rows;
		return ErrBadRequest // But I prefer to return ErrBadRequest;
	}

	if deadline.Before(time.Now()) {
		return ErrForbidden // After the deadline no voting;
	}

	// Step # 2: let's try to insert a new record into the 'voters' table;
	stmt = "INSERT INTO polls.voters (vote_id, user_id, created) VALUES(?, ?, toTimeStamp(now())) IF NOT EXISTS"
	m := make(map[string]interface{})
	applied, err := session.Query(stmt, vote_id, user_id).WithContext(ctx).MapScanCAS(m)
	if err != nil {
		return err
	}
	if !applied {
		return ErrForbidden // Looks like this voter has voted earlier;
	}

	// Step # 3: let's increment the 'co_count' for the specified contender in the 'votes' table.
	stmt = `UPDATE polls.votes SET co_count = co_count + 1, co_updated = toTimeStamp(now())
		WHERE vote_id = ? AND co_id = ? IF EXISTS`
	m = make(map[string]interface{})
	applied, err = session.Query(stmt, vote_id, co_id).WithContext(ctx).MapScanCAS(m)
	if !(err == nil && applied) {
		// If this failed, the voter has the right to vote again.
		// It means that the voter's 'user_id' must be removed from the 'voters' table.
		stmt = "DELETE FROM polls.voters WHERE vote_id = ? AND user_id = ?"
		session.Query(stmt, vote_id, user_id).WithContext(ctx).Exec()
		return err
	}

	// Success!
	return nil
}

//////////////////////
//
// GET SERVICE STATUS
//
//////////////////////

func (b *basicVoteService) GetServiceStatus(ctx context.Context) *HealthStatus {

	var hs HealthStatus

	// Let's check if the system has enough memory.
	mem := getAvailableMem()
	// 0 means that func failed to read the sys data, and we cannot rely on it.
	// fmt.Printf("HealthCheck, free mem = %d, LOW_MEM_THRESHOLD = %d",	mem, LOW_MEM_THRESHOLD)
	if mem > 0 && mem < LOW_MEM_THRESHOLD {
		hs.Status = http.StatusServiceUnavailable
		hs.Message = SERVICE_STATUS_LOW_MEMORY
		return &hs
	}

	// Let's check if the database is available.
	if b.db == nil {
		hs.Status = http.StatusServiceUnavailable
		hs.Message = SERVICE_STATUS_NO_DATABASE
		return &hs
	} else {
		// Let's try to connect to the database.
		session, err := b.db.CreateSession()
		if err != nil {
			hs.Status = http.StatusServiceUnavailable
			hs.Message = SERVICE_STATUS_NO_DATABASE
			return &hs
		}
		session.Close()
	}

	hs.Status = http.StatusOK
	hs.Message = SERVICE_STATUS_OK
	return &hs
}

//////////////////////////
//
// NEW BASIC VOTE SERVICE
//
//////////////////////////

// NewBasicVoteService returns a naive, stateless implementation of VoteService.
func NewBasicVoteService(db *gocql.ClusterConfig) VoteService {
	return &basicVoteService{
		db: db,
	}
}

////////////////////
//
// NEW VOTE SERVICE
//
////////////////////

// New returns a VoteService with all of the expected middleware wired in.
func New(db *gocql.ClusterConfig, middleware []Middleware) VoteService {
	var svc VoteService = NewBasicVoteService(db)
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}

/////////////////////
//
// GET AVAILABLE MEM
//
/////////////////////

// IMPORTANT! This is good for Linux only!

func getAvailableMem() uint64 {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}

	s := string(data)
	i := strings.Index(s, "MemAvailable")
	line := s[i:(i + 60)]
	lines := strings.Fields(line)
	memSizeAvailable, err := strconv.ParseUint(lines[1], 10, 64) // It returns uint64.
	if err != nil {
		return 0
	}

	return memSizeAvailable
}

// --- END OF FILE ---
