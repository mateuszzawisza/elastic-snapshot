package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mateuszzawisza/elastic-snapshot/snapshot"
)

var address = flag.String("address", "localhost:9200", "elasticsearch address:port")
var action = flag.String("action", "list", "(list|create)")
var repo = flag.String("repo", "my_repo", "snapshot repo")

func init() {
	if len(*repo) == 0 {
		log.Panicf("Repo not set")
	}
}

func main() {
	flag.Parse()
	switch *action {
	case "create":
		snap_name := fmt.Sprintf("snapshot_%s", time.Now)
		snapshot.CreateSnapshot(*address, *repo, snap_name)
	case "list":
		snapshot.ListSnapshots(*address, *repo)
	}
}
