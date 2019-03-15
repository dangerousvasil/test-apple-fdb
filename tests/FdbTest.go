package tests

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"io"
	"log"
	"os"
	"sync"
	"test-apple-fdb/program"
	"time"
)

func (t *FdbTest) timer() {
	var last int64
	for {
		log.Println("All", t.counter, "correct", t.correct, "In second ", (t.counter - last))
		last = t.counter
		time.Sleep(time.Second)
	}
}
func (t *FdbTest) timerMinute() {
	var last int64
	for {
		log.Println("[MINUTE] All", t.counter, "correct", t.correct, "In minute ", (t.counter - last))
		last = t.counter
		time.Sleep(60 * time.Second)
	}
}

type FdbTest struct {
	fdb         fdb.Database
	err         error
	wg          sync.WaitGroup
	titles      chan string
	SkipLine    int64
	LogFilename string
	FDBCluster  string
	counter     int64
	correct     int64
}

func (t *FdbTest) Init() {
	fdb.MustAPIVersion(600)
	if t.fdb, t.err = fdb.Open(t.FDBCluster, []byte(`DB`)); t.err != nil {
		panic(t.err)
	}
	t.titles = make(chan string, 2000)

}

func (t *FdbTest) TestLine(line string) error {
	t.counter++
	ttlpart := program.Title{}

	err := json.Unmarshal([]byte(line), &ttlpart)
	if err != nil {
		//log.Println(line)
		//log.Println(err)
		return err
	}
	//fmt.Println(ttlpart)
	aStr, err := t.fdb.Transact(func(tr fdb.Transaction) (result interface{}, err error) {

		trOpt := tr.Options()
		trOpt.SetReadLockAware()

		h1 := program.GetHkey1(ttlpart)
		titles, err := directory.CreateOrOpen(tr, []string{`titleFiz`}, nil)
		if err != nil {
			return
		}

		listHK1 := titles.Sub(`hkey1`)
		listId := titles.Sub(`list`)

		rr := tr.GetRange(listHK1.Sub(*h1), fdb.RangeOptions{})
		ri := rr.Iterator()
		ttlCnt := 0
		rezz := []string{}
		for ri.Advance() {
			var kv fdb.KeyValue
			kv, err = ri.Get()
			if err != nil {
				fmt.Printf("Unable to read next value: %v\n", err)
				return
			}
			ttlCnt++

			var unkey tuple.Tuple
			unkey, err = listHK1.Unpack(kv.Key)

			if err != nil {
				log.Println(err)
				return
			}

			keyId := listId.Pack(tuple.Tuple{unkey[1].(string)})

			fbyte := tr.Get(keyId)

			var data []byte
			data, err = fbyte.Get()
			if err != nil {
				return
			}

			rezz = append(rezz, string(data))
		}
		return rezz, err
	})

	if err != nil {
		return err
	}

	if len(aStr.([]string)) > 1 {
		return errors.New(fmt.Sprintf("Panic double ttitle %v \r\n %v", line, aStr))
	} else if len(aStr.([]string)) == 0 {
		log.Println(line)
		log.Println("Not found")
		log.Fatalln("Test not passed")
	}

	if string(aStr.([]string)[0]) == line {
		t.correct++
	} else {
		log.Println(line)
		log.Fatalln("Test not passed")
	}
	return nil
}

func (t *FdbTest) readFile() {
	defer close(t.titles)

	file, err := os.Open(t.LogFilename)
	defer file.Close()

	if err != nil {
		log.Fatalln(err)
	}

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)
	var i int64
	for {
		var buffer bytes.Buffer

		var l []byte
		var isPrefix bool
		for {
			l, isPrefix, err = reader.ReadLine()
			buffer.Write(l)

			// If we've reached the end of the line, stop reading.
			if !isPrefix {
				break
			}

			// If we're just at the EOF, break
			if err != nil {
				break
			}
		}

		if err == io.EOF {
			break
		}

		line := buffer.String()

		if t.SkipLine < i {
			t.titles <- line
		}
		i++
	}

	if err != io.EOF {
		fmt.Printf(" > Failed!: %v\n", err)
	}

	return
}
func (t *FdbTest) worker() {
	defer t.wg.Done()

	for {
		select {
		case s, ok := <-t.titles:
			if ok {
				t.TestLine(s)
			} else {
				return
			}
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
