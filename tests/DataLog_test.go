package tests

import (
	"log"
	"runtime"
	"testing"
)

// if have version problem run
//  go get github.com/apple/foundationdb/bindings/go@release-6.0
func TestDataLog(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	tester := FdbTest{
		SkipLine:    0,
		LogFilename: `/u03/iron_logs/iron_hdd_000_01.log`,
		FDBCluster:  `/etc/foundationdb/fdb.cluster`,
	}
	tester.Init()

	go tester.readFile()

	for n := 0; n < runtime.NumCPU()*4; n++ {
		tester.wg.Add(1)
		go tester.worker()
	}

	go tester.timer()
	go tester.timerMinute()

	tester.wg.Wait()
}
