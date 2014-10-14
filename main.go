package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mateuszzawisza/elastic-snapshot/snapshot"
)

const Version = "0.0.3"

var address = flag.String("address", "http://localhost:9200", "elasticsearch address:port")
var action = flag.String("action", "action", "(list|create)")
var repo = flag.String("repo", "my_repo", "snapshot repo")
var version = flag.Bool("version", false, "Print version and exit")

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
	case "create":
		snap_name := fmt.Sprintf("snapshot_%d", time.Now().Unix())
		snapshot.CreateSnapshot(*address, *repo, snap_name)
	case "restore":
		err := snapshot.RestoreLastSnapshot(*address, *repo)
		if err != nil {
			log.Fatal("Restore failed")
		}
	case "list":
		snapshots := snapshot.ListSnapshots(*address, *repo).Snapshots
		for _, snapshot := range snapshots {
			fmt.Println(snapshot.Snapshot)
		}
	}
}
