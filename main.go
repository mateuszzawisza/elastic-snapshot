package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mateuszzawisza/elastic-brain-surgeon/clusterstatus"
	"github.com/mateuszzawisza/elastic-snapshot/snapshot"
)

const Version = "0.0.6"

var address = flag.String("address", "http://localhost:9200", "elasticsearch address:port")
var action = flag.String("action", "action", "(create-repo|list|create|restore|clean-old)")
var keep = flag.Int("keep-snapshots", 720, "How many snapshots to keep")
var repo = flag.String("repo", "my_repo", "snapshot repo")
var bucketName = flag.String("bucket-name", "bucket", "Bucket name for the repository")
var basePath = flag.String("base-path", "", "path in bucket")
var version = flag.Bool("version", false, "Print version and exit")
var masterOnly = flag.Bool("master-only", false, "Perform action only if current node is master node")

func init() {
	if len(*repo) == 0 {
		log.Panicf("Repo not set")
	}
}

func main() {
	flag.Parse()
	if *version {
		fmt.Printf("Version: %s\n", Version)
		return
	}
	switch *action {
	default:
		return
	case "create-repo":
		repoExists, err := snapshot.CheckRepo(*address, *repo)
		if err != nil {
			log.Fatal("Could not verify if repo exists. Got: %s", err)
		}
		if repoExists {
			log.Printf("Repo exists.")
		} else {
			snapshot.CreateRepo(*address, *repo, *bucketName, *basePath)
		}
	case "create":
		if *masterOnly {
			if amI := clusterstatus.AmIMaster(*address); amI == false {
				log.Println("I'm not a master. Exiting.")
				return
			}
		}
		snap_name := fmt.Sprintf("snapshot_%d", time.Now().Unix())
		snapshot.CreateSnapshot(*address, *repo, snap_name)
	case "restore":
		err := snapshot.RestoreLastSnapshot(*address, *repo)
		if err != nil {
			log.Fatal("Restore failed")
		}
	case "clean-old":
		if *masterOnly {
			if amI := clusterstatus.AmIMaster(*address); amI == false {
				log.Println("I'm not a master. Exiting.")
				return
			}
		}
		err := snapshot.SnapshotRetention(*address, *repo, *keep)
		if err != nil {
			log.Fatal("Cleanup failed. Got: %s", err)
		}
	case "list":
		snapshotList, err := snapshot.ListSnapshots(*address, *repo)
		if err != nil {
			log.Fatal("Failed to list snapshots. Error: %s", err)
		}
		snapshots := snapshotList.Snapshots
		for _, snapshot := range snapshots {
			fmt.Println(snapshot.Snapshot)
		}
	}
}
