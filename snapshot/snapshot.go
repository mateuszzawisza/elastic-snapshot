package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type SnapshotRequest struct {
	uri          string
	requestPath  string
	method       string
	pathSettings map[string]string
	data         string
}

type listSnapshotsJSON struct {
	Snapshots []struct {
		Snapshot          string   `json:"snapshot"`
		Indices           []string `json:"indices"`
		State             []string `json:"state"`
		StartTime         string   `json:"start_time"`
		StartTimeInMillis int      `json:"start_time_in_millis"`
		EndTime           string   `json:"end_time"`
		EndTimeInMillis   int      `json:"end_time_in_millis"`
		DurationInMillis  int      `json:"DurationInMillis"`
		Failurs           []string `json:"failures"`
		Shards            struct {
			Total      int `json:"total"`
			Failed     int `json:"failed"`
			Successful int `json:"successful"`
		} `json: "shards"`
	} `json:"snapshots"`
}

var CheckRepoRequest SnapshotRequest = SnapshotRequest{
	"localhost:9200",
	"_snapshot/{{repo_name}}",
	"GET",
	map[string]string{},
	"",
}

var CreateRepoRequest SnapshotRequest = SnapshotRequest{
	"localhost:9200",
	"_snapshot/{{repo_name}}",
	"PUT",
	map[string]string{},
	`{
    "type": "s3",
    "settings": {
        "bucket": "{{bucket_name}}",
        "base_path": "{{base_path}}",
        "region": "{{region}}"
    }
}`,
}

var CreateSnapshotRequest SnapshotRequest = SnapshotRequest{
	"localhost:9200",
	"_snapshot/{{repo_name}}/{{snapshot_name}}?wait_for_completion=true",
	"PUT",
	map[string]string{},
	"",
}

var RestoreSnapshotRequest SnapshotRequest = SnapshotRequest{
	"localhost:9200",
	"_snapshot/{{repo_name}}/{{snapshot_name}}/_restore",
	"POST",
	map[string]string{},
	"",
}

var ListSnapshotsRequest SnapshotRequest = SnapshotRequest{
	"localhost:9200",
	"_snapshot/{{repo_name}}/_all",
	"GET",
	map[string]string{},
	"",
}

var DeleteSnapshotsRequest SnapshotRequest = SnapshotRequest{
	"localhost:9200",
	"_snapshot/{{repo_name}}/{{snapshot_name}}",
	"DELETE",
	map[string]string{},
	"",
}

var requestFailed = errors.New("Request failed")

func (r *SnapshotRequest) setPath() {
	path := r.requestPath
	for name, value := range r.pathSettings {
		nameMark := fmt.Sprintf("{{%s}}", name)
		path = strings.Replace(path, nameMark, value, 1)
	}
	r.requestPath = path
}

func (r *SnapshotRequest) setData() {
	data := r.data
	for name, value := range r.pathSettings {
		nameMark := fmt.Sprintf("{{%s}}", name)
		data = strings.Replace(data, nameMark, value, 1)
	}
	r.data = data
}

func (r *SnapshotRequest) perform() (*http.Response, error) {
	r.setPath()
	r.setData()
	client := &http.Client{}
	requestURL := fmt.Sprintf("%s/%s", r.uri, r.requestPath)
	body := strings.NewReader(r.data)
	req, err := http.NewRequest(r.method, requestURL, body)
	if err != nil {
		return nil, err
	}
	response, connectionErr := client.Do(req)
	if connectionErr != nil {
		return nil, connectionErr
	}
	if 500 <= response.StatusCode && response.StatusCode < 600 {
		return response, requestFailed
	}
	return response, nil
}

func CheckRepo(url, repoName string) (bool, error) {
	request := CheckRepoRequest
	request.uri = url
	request.pathSettings["repo_name"] = repoName
	response, err := request.perform()
	if err != nil {
		return false, err
	}
	switch response.StatusCode {
	default:
		log.Printf("Got status: %s - %d", response.StatusCode, response.Status)
		return false, nil
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	}
}

func CreateRepo(url, repoName, bucketName, region, basePath string) {
	request := CreateRepoRequest
	request.uri = url
	request.pathSettings["repo_name"] = repoName
	request.pathSettings["bucket_name"] = bucketName
	request.pathSettings["base_path"] = basePath
	request.pathSettings["region"] = region
	response, err := request.perform()
	defer response.Body.Close()
	if err != nil {
		body, _ := ioutil.ReadAll(response.Body)
		log.Panicf("Failed on create repo request. Err: %v, %s", err, string(body))
	}
}

func CreateSnapshot(url, repoName, snapName string) {
	request := CreateSnapshotRequest
	request.uri = url
	request.pathSettings["repo_name"] = repoName
	request.pathSettings["snapshot_name"] = snapName
	response, err := request.perform()
	defer response.Body.Close()
	if err != nil {
		body, _ := ioutil.ReadAll(response.Body)
		log.Panicf("Failed on create snapshot request. Err: %v, %s", err, string(body))
	}
}

func RestoreSnapshot(url, repoName, snapName string) {
	request := RestoreSnapshotRequest
	request.uri = url
	request.pathSettings["repo_name"] = repoName
	request.pathSettings["snapshot_name"] = snapName
	response, err := request.perform()
	defer response.Body.Close()
	if err != nil {
		body, _ := ioutil.ReadAll(response.Body)
		log.Panicf("Failed on restore snapshot request. Err: %v, %s", err, string(body))
	}
}

func RestoreLastSnapshot(url, repoName string) error {
	snapshots, err := ListSnapshots(url, repoName)
	if err != nil {
		return err
	}
	lastSnapshot, errSnapshot := findLastSnapshot(snapshots)
	if errSnapshot != nil {
		return errSnapshot
	}
	RestoreSnapshot(url, repoName, lastSnapshot)
	return nil
}

func DeleteSnapshot(url, repoName, snapName string) {
	request := DeleteSnapshotsRequest
	request.uri = url
	request.pathSettings["repo_name"] = repoName
	request.pathSettings["snapshot_name"] = snapName
	response, err := request.perform()
	defer response.Body.Close()
	if err != nil {
		body, _ := ioutil.ReadAll(response.Body)
		log.Panicf("Failed on delete snapshot request. Err: %v, %s", err, string(body))
	}
}

func SnapshotRetention(url, repoName string, snapshotsToKeep int) error {
	snapshotsList, err := ListSnapshots(url, repoName)
	if err != nil {
		return err
	}
	snapshots := snapshotsList.Snapshots
	if snapshotsToDeleteCount := (len(snapshots) - snapshotsToKeep); snapshotsToDeleteCount < 0 {
		return nil
	} else {
		snapshotsToDelete := snapshots[0:snapshotsToDeleteCount]
		for _, snapshot := range snapshotsToDelete {
			DeleteSnapshot(url, repoName, snapshot.Snapshot)
		}
		return nil
	}
}

func ListSnapshots(url, repoName string) (listSnapshotsJSON, error) {
	request := ListSnapshotsRequest
	request.uri = url
	request.pathSettings["repo_name"] = repoName
	response, err := request.perform()
	if err != nil {
		return listSnapshotsJSON{}, err
	}
	js := parseListSnapshotsResponse(response)
	return js, nil
}

func parseListSnapshotsResponse(response *http.Response) listSnapshotsJSON {
	var js listSnapshotsJSON
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Error reading list snapshot response")
	}
	json.Unmarshal(body, &js)
	return js
}

func findLastSnapshot(snapshots listSnapshotsJSON) (string, error) {
	snapshotsCount := len(snapshots.Snapshots)
	if snapshotsCount > 0 {
		return snapshots.Snapshots[snapshotsCount-1].Snapshot, nil
	} else {
		return "", errors.New("Last snapshot could not be found")
	}

}
