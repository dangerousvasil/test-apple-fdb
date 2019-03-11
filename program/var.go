package program

import (
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
