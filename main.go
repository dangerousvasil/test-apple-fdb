package main

import (
	"flag"
	"log"
	"test-apple-fdb/program"
)

var (
	clusterFile = flag.String("c", `/etc/foundationdb/fdb_replica.cluster`, "Using cluster file")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	log.Println("Using cluster file", *clusterFile)

	p := new(program.Program)
	p.Run(*clusterFile)
}
