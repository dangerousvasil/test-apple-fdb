package program

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"runtime"
	"sync"
	"time"
)

type Title struct {
	FileId     int64
	PartnerId  int64
	Unit       int64
	LastName   string
	FirstName  string
	MiddleName string
	MaidenName string
	Birthday   time.Time
	BirthPlace string
	DocType    int
	DocNo      string
	DocPlace   string
	DocDate    time.Time
	DocEndDate time.Time
	Sex        int
	Address    addrs
	Id         string
}

type addrs struct {
	Reg  *addr
	Fact *addr
}

type addr struct {
	ZipCode   string
	Country   string
	Region    string
	City      string
	District  string
	Statement string
	Street    string
	House     string
	Block     string
	Build     string
	Flat      string
}

type Program struct {
	Buf chan Title
	fdb fdb.Database
	wWg sync.WaitGroup
	err error
	run bool
}

func (p *Program) Run(file string) {
	p.run = true
	fdb.MustAPIVersion(600)
	if p.fdb, p.err = fdb.Open(file, []byte(`DB`)); p.err != nil {
		panic(p.err)
	}
	p.Buf = make(chan Title, 10000)
	go p.SetBuf()
	//p.wWg.Add(1)
	//go p.worker()
	for i := 0; i < runtime.NumCPU()*4; i++ {
		p.wWg.Add(1)
		go p.worker()
	}
	p.wWg.Wait()
}

func (p *Program) SetBuf() {
	for p.run {
		p.Buf <- p.GenerateTitle()
	}
	close(p.Buf)
}

func (p *Program) worker() {
	defer func() {
		p.wWg.Done()
	}()
	for {
		select {
		case t, ok := <-p.Buf:
			if ok {
				p.Save(t)
			} else {
				return
			}
			break
		default:
			time.Sleep(1)
		}
	}
}
