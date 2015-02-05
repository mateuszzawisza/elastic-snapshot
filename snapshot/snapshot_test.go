package snapshot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

func TestCheckRepoFound(t *testing.T) {
	repoName := "super_cluster_repository"
	const expectedURI = "/_snapshot/super_cluster_repository"
	const expectedHTTPMethod = "GET"
	var receivedURI string
	var receivedHTTPMethod string

	es := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURI = r.RequestURI
		receivedHTTPMethod = r.Method
		fmt.Fprintln(w, "`{}`")
	}))
	defer es.Close()
	response, err := CheckRepo(es.URL, repoName)
	if err != nil {
		t.Fatalf("Received unexpected error. Got %s", err)
	}
	if receivedURI != expectedURI {
		t.Fatalf("Request URI not matched. Got %s. Expected: %s", receivedURI, expectedURI)
	}
	if !response {
		t.Fatalf("Repo exists but not found!")
	}
}

func TestCheckRepoNotFound(t *testing.T) {
	repoName := "super_cluster_repository"
	const expectedURI = "/_snapshot/super_cluster_repository"
	const expectedHTTPMethod = "GET"
	var receivedURI string
	var receivedHTTPMethod string

	es := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURI = r.RequestURI
		receivedHTTPMethod = r.Method
		http.NotFound(w, r)
	}))
	defer es.Close()
	response, err := CheckRepo(es.URL, repoName)
	if err != nil {
		t.Fatalf("Received unexpected error. Got %s", err)
	}
	if receivedURI != expectedURI {
		t.Fatalf("Request URI not matched. Got %s. Expected: %s", receivedURI, expectedURI)
	}
	if response {
		t.Fatalf("Repo doesn't exists but was found!")
	}
}

func TestCreateRepo(t *testing.T) {
	repoName := "super_cluster_repository"
	bucketName := "elasticsearch-europe"
  region := "us-west-1"
	basePath := "super_cluster"
	const expectedURI = "/_snapshot/super_cluster_repository"
	const expectedHTTPMethod = "PUT"
	const expectedData = `{
    "type": "s3",
    "settings": {
        "bucket": "elasticsearch-europe",
        "base_path": "super_cluster",
        "region": "us-west-1"
    }
}`
	var receivedURI string
	var receivedHTTPMethod string
	var receivedData []byte

	es := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURI = r.RequestURI
		receivedHTTPMethod = r.Method
		receivedData, _ = ioutil.ReadAll(r.Body)
		fmt.Fprintln(w, "`{}`")
	}))
	defer es.Close()
	CreateRepo(es.URL, repoName, bucketName, region, basePath)
	if receivedURI != expectedURI {
		t.Fatalf("Request URI not matched. Got %s. Expected: %s", receivedURI, expectedURI)
	}
	if receivedHTTPMethod != expectedHTTPMethod {
		t.Fatalf("Request Method not matched. Got %s. Expected: %s", receivedHTTPMethod, expectedHTTPMethod)
	}
	if string(receivedData) != expectedData {
		t.Fatalf("Request data not matched. Got %s. Expected: %s", receivedData, expectedData)
	}
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
	listSnapshotResponse := loadJSONFromFile("list_snapshot_response_test.json")

	const expectedURI = "/_snapshot/test_repo/snapshot_1414576801/_restore"
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
	listSnapshotResponse := loadJSONFromFile("list_snapshot_response_test.json")
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

	snapshots, _ := ListSnapshots(ts.URL, "test_repo")
	if receivedURI != expectedURI {
		t.Fatalf("Request URI not matched. Got %s. Expected: %s", receivedURI, expectedURI)
	}
	if receivedHTTPMethod != expectedHTTPMethod {
		t.Fatalf("Request Method not matched. Got %s. Expected: %s", receivedHTTPMethod, expectedHTTPMethod)
	}
	if snapshots.Snapshots[0].Snapshot != "snapshot_1414514679" {
		t.Fatalf("Snaphsot name mismatch")
	}
	if snapshots.Snapshots[len(snapshots.Snapshots)-1].Snapshot != "snapshot_1414576801" {
		t.Fatalf("Snaphsot name mismatch")
	}
}

func TestDeleteSnapshot(t *testing.T) {
	const expectedURI = "/_snapshot/test_repo/snap_1412944813"
	const expectedHTTPMethod = "DELETE"
	var receivedURI string
	var receivedHTTPMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURI = r.RequestURI
		receivedHTTPMethod = r.Method
		fmt.Fprintln(w, "`{}`")
	}))
	defer ts.Close()

	DeleteSnapshot(ts.URL, "test_repo", "snap_1412944813")
	if receivedURI != expectedURI {
		t.Fatalf("Request URI not matched. Got %s. Expected: %s", receivedURI, expectedURI)
	}
	if receivedHTTPMethod != expectedHTTPMethod {
		t.Fatalf("Request Method not matched. Got %s. Expected: %s", receivedHTTPMethod, expectedHTTPMethod)
	}
}

func TestSnapshotRetention(t *testing.T) {
	listSnapshotResponse := loadJSONFromFile("list_snapshot_response_test.json")
	const expectedDeletes = 8
	receivedDeletes := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHTTPMethod := r.Method
		if receivedHTTPMethod == "GET" {
			fmt.Fprintln(w, listSnapshotResponse)
		} else if receivedHTTPMethod == "DELETE" {
			receivedDeletes = receivedDeletes + 1
			fmt.Fprintln(w, "`{}`")
		}
	}))
	defer ts.Close()

	SnapshotRetention(ts.URL, "test_repo", 10)
	if expectedDeletes != receivedDeletes {
		t.Fatalf("Expected to receive %d deletes. Got: %d", expectedDeletes, receivedDeletes)
	}
}

func TestSnapshotRetentionWithConnectionError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHTTPMethod := r.Method
		if receivedHTTPMethod == "GET" {
		} else if receivedHTTPMethod == "DELETE" {
			fmt.Fprintln(w, "`{}`")
		}
	}))
	ts.Close() // we're closing connection immediately to check how code will react

	if err := SnapshotRetention(ts.URL, "test_repo", 10); err == nil {
		t.Fatalf("Expected error but here weren't any")
	}
}

func TestSnapshotRetentionWithError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHTTPMethod := r.Method
		if receivedHTTPMethod == "GET" {
			http.Error(w, "Server Error", http.StatusInternalServerError)
		} else if receivedHTTPMethod == "DELETE" {
			fmt.Fprintln(w, "`{}`")
		}
	}))
	defer ts.Close()

	if err := SnapshotRetention(ts.URL, "test_repo", 10); err == nil {
		t.Fatalf("Expected error but here weren't any")
	}
}

func loadJSONFromFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	stat, _ := file.Stat()
	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
