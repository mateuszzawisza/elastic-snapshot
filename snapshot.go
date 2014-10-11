package snapshot

import (
	"fmt"
	"strings"
)

type SnapshotRequest struct {
	host         string
	port         int
	requestPath  string
	method       string
	pathSettings map[string]string
}

var CreateSnapshotRepoRequest SnapshotRequest = SnapshotRequest{
	"localhost",
	9200,
	"_snapshot/{{repo_name}}/{{snapshot_name}}",
	"PUT",
	map[string]string{},
}

func (r *SnapshotRequest) setPath() {
	path := r.requestPath
	for name, value := range r.pathSettings {
		nameMark := fmt.Sprintf("{{%s}}", name)
		path = strings.Replace(path, nameMark, value, 1)
	}
	r.requestPath = path
}
