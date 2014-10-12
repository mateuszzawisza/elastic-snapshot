package snapshot

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type SnapshotRequest struct {
	uri          string
	requestPath  string
	method       string
	pathSettings map[string]string
}

var CreateSnapshotRequest SnapshotRequest = SnapshotRequest{
	"localhost:9200",
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

func (r *SnapshotRequest) perform() {
	r.setPath()
	client := &http.Client{}
	requestURL := fmt.Sprintf("%s/%s", r.uri, r.requestPath)
	req, err := http.NewRequest(r.method, requestURL, nil)
	if err != nil {
		log.Panic("Error creating request object")
	}
	_, connectionErr := client.Do(req)
	if connectionErr != nil {
		log.Panic("Error connecting with Elasticsearch")
	}
}

func createSnapshot(url, repoName, snapName string) {
	request := CreateSnapshotRequest
	request.uri = url
	request.pathSettings["repo_name"] = repoName
	request.pathSettings["snapshot_name"] = snapName
	request.perform()
}
