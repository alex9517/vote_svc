//  Created : 2024-Mar-14
// Modified : 2024-Mar-30

package service

import (
	// "bytes"
	// "context"
	// "fmt"
	"context"
	"fmt"
	"net/http"
	"time"

	// "os"
	"testing"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	// "github.com/hailocab/go-hostpool"
)

const ENTRY_NODE_1 = "172.16.70.31"

// const ENTRY_NODE_2 = "172.16.70.32"
const KEYSPACE_NAME = "polls"
const TESTDATA_DIR = "../../testdata/"

func connect_db() *gocql.ClusterConfig {
	cluster := gocql.NewCluster(ENTRY_NODE_1)
	cluster.Consistency = gocql.Quorum
	// cluster.ProtoVersion = 4
	cluster.Timeout = time.Second * 10
	cluster.ConnectTimeout = time.Second * 10
	// cluster.Keyspace = KEYSPACE_NAME
	// cluster.PoolConfig.HostSelectionPolicy = gocql.HostPoolHostPolicy(hostpool.New(nil))

	// cluster.Authenticator = gocql.PasswordAuthenticator{
	// Username: "user",
	// Password: "password",
	// }

	return cluster
}

func TestGetVoteData(t *testing.T) {
	testinfo := "test GetVoteData"
	// var vote_id int = 1
	var resources string = "https://ws4/votes/1/images/"

	cluster := connect_db()

	if cluster == nil {
		t.Errorf("cannot connect to database")
		return
	}

	svc := New(cluster, []Middleware{})

	t.Run(testinfo, func(t *testing.T) {
		// Case 1: vote_id = 1 (good), func must return 'VoteData' and err must be 'nil';
		res, err := svc.GetVoteData(context.Background(), 1)
		if err != nil {
			t.Errorf("test %v (case # 1) failed, error: %v", testinfo, err)
			return
		}

		if res.Resources != resources {
			t.Errorf("test %v (case # 1) failed, res.Resources is %v, but it should be %v", testinfo, res.Resources, resources)
			return
		}

		// Case 2: vote_id = 2 (bad), func must return 'nil' err must be 'ErrNotFound';
		_, err = svc.GetVoteData(context.Background(), 2)
		if err != ErrNotFound {
			t.Errorf("test %v (case # 2) failed, error: %v (must be %v)", testinfo, err, ErrNotFound)
			return
		}
	})
}

func TestAddVoteOk(t *testing.T) {
	testinfo := "test AddVoteOk"
	vote_id := 1
	co_id := 1
	co_count := 0
	random_user_id := uuid.NewString()
	cluster := connect_db()

	if cluster == nil {
		t.Errorf("cannot connect to database")
		return
	}

	svc := New(cluster, []Middleware{})

	session, err := cluster.CreateSession()
	if err != nil {
		t.Errorf("test %v failed, cannot create session, error: %v", testinfo, err)
		return
	}
	defer session.Close()

	// Get current 'co_count' related to 'vote_id' and 'co_id'; save it as 'old_co_count';
	stmt := "SELECT vote_id, co_id, co_count FROM polls.votes WHERE vote_id = ? AND co_id = ? LIMIT 1"

	err = session.Query(stmt, vote_id, co_id).WithContext(context.Background()).Consistency(gocql.One).Scan(&vote_id, &co_id, &co_count)
	if err != nil {
		t.Errorf("test %v failed, cannot get current co_count, error: %v", testinfo, err)
		return
	}
	old_co_count := co_count

	t.Run(testinfo, func(t *testing.T) {
		// Add one vote to 'co_id' related to 'vote_id';
		err := svc.UpdateVoteResults(context.Background(), vote_id, int16(co_id), random_user_id)
		if err != nil {
			t.Errorf("test %v failed, error: %v", testinfo, err)
			return
		}

		// Once again get current 'co_count' related to 'vote_id' and 'co_id'; it must be increased by one;
		stmt := "SELECT vote_id, co_id, co_count FROM polls.votes WHERE vote_id = ? AND co_id = ? LIMIT 1"
		err = session.Query(stmt, vote_id, co_id).WithContext(context.Background()).Consistency(gocql.One).Scan(&vote_id, &co_id, &co_count)
		if err != nil {
			t.Errorf("test %v failed, cannot get new co_count, error: %v", testinfo, err)
			return
		}

		if co_count != (old_co_count + 1) {
			t.Errorf("test %v failed, co_count = %v, but it should be %v", testinfo, co_count, (old_co_count + 1))
			return
		}
	})
}

func TestAddVoteReject(t *testing.T) {
	// This is similar to 'TestAddVoteOk', but we try to vote twice with the same user_id;
	// This should be impossible;
	testinfo := "test AddVoteReject"
	vote_id := 1
	co_id := 1
	co_count := 0

	// This is supposed to be already present in the 'polls.voters' table;
	user_id := "8b820e4e-8f43-4cd1-a8b3-f90c44e13ece"

	cluster := connect_db()

	if cluster == nil {
		t.Errorf("cannot connect to database")
		return
	}

	svc := New(cluster, []Middleware{})

	session, err := cluster.CreateSession()
	if err != nil {
		t.Errorf("test %v failed, cannot create session, error: %v", testinfo, err)
		return
	}
	defer session.Close()

	// Get current 'co_count' related to 'vote_id' and 'co_id'; save it as 'old_co_count';
	stmt := "SELECT vote_id, co_id, co_count FROM polls.votes WHERE vote_id = ? AND co_id = ? LIMIT 1"

	err = session.Query(stmt, vote_id, co_id).WithContext(context.Background()).Consistency(gocql.One).Scan(&vote_id, &co_id, &co_count)
	if err != nil {
		t.Errorf("test %v failed, cannot get current co_count, error: %v", testinfo, err)
		return
	}
	old_co_count := co_count

	t.Run(testinfo, func(t *testing.T) {
		// Case 1: Try to add one vote to 'co_id' related to 'vote_id';
		err := svc.UpdateVoteResults(context.Background(), vote_id, int16(co_id), user_id)
		if err != ErrForbidden {
			t.Errorf("test %v (case # 1) failed, error: %v (must be %v)", testinfo, err, ErrForbidden)
		}

		// Once again get current 'co_count' related to 'vote_id' and 'co_id'; it must be unchanged;
		stmt := "SELECT vote_id, co_id, co_count FROM polls.votes WHERE vote_id = ? AND co_id = ? LIMIT 1"
		err = session.Query(stmt, vote_id, co_id).WithContext(context.Background()).Consistency(gocql.One).Scan(&vote_id, &co_id, &co_count)
		if err != nil {
			t.Errorf("test %v (case # 1) failed, cannot get new co_count, error: %v", testinfo, err)
			return
		}

		if co_count == (old_co_count + 1) {
			t.Errorf("test %v failed, co_count = %v, but it should be %v", testinfo, co_count, old_co_count)
			return
		}

		// Case 2: Try to add one vote to 'co_id' related to non-existent 'vote_id';

	})
}

func TestHealthCheck(t *testing.T) {
	testinfo := "test HealthCheck"
	cluster := connect_db()

	svc := New(cluster, []Middleware{})

	t.Run(testinfo, func(t *testing.T) {
		res := svc.GetServiceStatus(context.Background())
		if res.Status != http.StatusOK {
			t.Errorf("test %v failed, res is %v", testinfo, res.Status)
		}
	})
}

func TestGetFreeMem(t *testing.T) {
	testinfo := "test getFreeMeme"

	t.Run(testinfo, func(t *testing.T) {
		mem := getAvailableMem()
		fmt.Printf("Free memory: %v", mem)
		if mem == 0 {
			t.Errorf("test %v failed, mem is %v", testinfo, mem)
		}
	})
}

// --- END OF FILE ---
