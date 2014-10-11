package snapshot

import "testing"

func TestInterpolateRequestPath(t *testing.T) {
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
