package snapshot

import (
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestCreateSnapshot(t *testing.T) {
	const expectedURI = "/_snapshot/test_repo/snap_1"
	const expectedHTTPMethod = "PUT"
	var receivedURI string
	var receivedHTTPMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURI = r.RequestURI
		receivedHTTPMethod = r.Method
		fmt.Fprintln(w, "{}")
	}))
	defer ts.Close()

	createSnapshot(ts.URL, "test_repo", "snap_1")
	if receivedURI != expectedURI {
		t.Fatalf("Request URI not matched. Got %s. Expected: %s", receivedURI, expectedURI)
	}
	if receivedHTTPMethod != expectedHTTPMethod {
		t.Fatalf("Request Method not matched. Got %s. Expected: %s", receivedHTTPMethod, expectedHTTPMethod)
	}
}

func TestCreateRepo(t *testing.T) {
}

func TestDeleteSnapshot(t *testing.T) {
}

func TestListSnapshots(t *testing.T) {
}
