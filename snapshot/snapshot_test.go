package snapshot

import (
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"
)

const listSnapshotResponse string = `{
  "snapshots" : [ {
    "snapshot" : "snap_1412944591",
    "indices" : [ "twitter" ],
    "state" : "SUCCESS",
    "start_time" : "2014-10-10T12:36:31.531Z",
    "start_time_in_millis" : 1412944591531,
    "end_time" : "2014-10-10T12:36:33.115Z",
    "end_time_in_millis" : 1412944593115,
    "duration_in_millis" : 1584,
    "failures" : [ ],
    "shards" : {
      "total" : 5,
      "failed" : 0,
      "successful" : 5
    }
  }, {
    "snapshot" : "snap_1412944813",
    "indices" : [ "twitter" ],
    "state" : "SUCCESS",
    "start_time" : "2014-10-10T12:40:13.904Z",
    "start_time_in_millis" : 1412944813904,
    "end_time" : "2014-10-10T12:40:14.730Z",
    "end_time_in_millis" : 1412944814730,
    "duration_in_millis" : 826,
    "failures" : [ ],
    "shards" : {
      "total" : 5,
      "failed" : 0,
      "successful" : 5
    }
  }]
}`

func TestSetPath(t *testing.T) {
	const expectedRequestPath = "_snapshot/test_repo/snap_1"
	snapReqTest := new(SnapshotRequest)
	snapReqTest.requestPath = "_snapshot/{{repo_name}}/{{snapshot_name}}"
	snapReqTest.pathSettings = make(map[string]string)
	snapReqTest.pathSettings["repo_name"] = "test_repo"
	snapReqTest.pathSettings["snapshot_name"] = "snap_1"
	snapReqTest.setPath()
	interpolatedPath := snapReqTest.requestPath
	if interpolatedPath != expectedRequestPath {
		t.Fatalf("Request path not set properly. Got %s. Expected: %s", interpolatedPath, expectedRequestPath)
	}
}

func TestCreateRepo(t *testing.T) {
}

func TestCreateSnapshot(t *testing.T) {
	const expectedURI = "/_snapshot/test_repo/snap_1?wait_for_completion=true"
	const expectedHTTPMethod = "PUT"
	var receivedURI string
	var receivedHTTPMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURI = r.RequestURI
		receivedHTTPMethod = r.Method
		fmt.Fprintln(w, "`{}`")
	}))
	defer ts.Close()

	CreateSnapshot(ts.URL, "test_repo", "snap_1")
	if receivedURI != expectedURI {
		t.Fatalf("Request URI not matched. Got %s. Expected: %s", receivedURI, expectedURI)
	}
	if receivedHTTPMethod != expectedHTTPMethod {
		t.Fatalf("Request Method not matched. Got %s. Expected: %s", receivedHTTPMethod, expectedHTTPMethod)
	}
}

func TestRestoreSnapshot(t *testing.T) {
	const expectedURI = "/_snapshot/test_repo/snap_1/_restore"
	const expectedHTTPMethod = "POST"
	var receivedURI string
	var receivedHTTPMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURI = r.RequestURI
		receivedHTTPMethod = r.Method
		fmt.Fprintln(w, "`{}`")
	}))
	defer ts.Close()

	RestoreSnapshot(ts.URL, "test_repo", "snap_1")
	if receivedURI != expectedURI {
		t.Fatalf("Request URI not matched. Got %s. Expected: %s", receivedURI, expectedURI)
	}
	if receivedHTTPMethod != expectedHTTPMethod {
		t.Fatalf("Request Method not matched. Got %s. Expected: %s", receivedHTTPMethod, expectedHTTPMethod)
	}
}

func TestRestoreLastSnapshot(t *testing.T) {
	const expectedURI = "/_snapshot/test_repo/snap_1412944813/_restore"
	const expectedHTTPMethod = "POST"
	var receivedURI string
	var receivedHTTPMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURI = r.RequestURI
		receivedHTTPMethod = r.Method

		if receivedHTTPMethod == "GET" {
			fmt.Fprintln(w, listSnapshotResponse)
		} else {
			fmt.Fprintln(w, "`{}`")
		}
	}))
	defer ts.Close()

	RestoreLastSnapshot(ts.URL, "test_repo")
	if receivedURI != expectedURI {
		t.Fatalf("Request URI not matched. Got %s. Expected: %s", receivedURI, expectedURI)
	}
	if receivedHTTPMethod != expectedHTTPMethod {
		t.Fatalf("Request Method not matched. Got %s. Expected: %s", receivedHTTPMethod, expectedHTTPMethod)
	}
}

func TestListSnapshots(t *testing.T) {
	const expectedURI = "/_snapshot/test_repo/_all"
	const expectedHTTPMethod = "GET"
	var receivedURI string
	var receivedHTTPMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURI = r.RequestURI
		receivedHTTPMethod = r.Method
		fmt.Fprintln(w, listSnapshotResponse)
	}))
	defer ts.Close()

	snapshots := ListSnapshots(ts.URL, "test_repo")
	if receivedURI != expectedURI {
		t.Fatalf("Request URI not matched. Got %s. Expected: %s", receivedURI, expectedURI)
	}
	if receivedHTTPMethod != expectedHTTPMethod {
		t.Fatalf("Request Method not matched. Got %s. Expected: %s", receivedHTTPMethod, expectedHTTPMethod)
	}
	if snapshots.Snapshots[0].Snapshot != "snap_1412944591" {
		t.Fatalf("Snaphsot name mismatch")
	}
	if snapshots.Snapshots[1].Snapshot != "snap_1412944813" {
		t.Fatalf("Snaphsot name mismatch")
	}
}

func TestDeleteSnapshot(t *testing.T) {
}
