package program

import (
	"encoding/json"
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"log"
)

func (p *Program) Save(t Title) {

	body, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		return
	}

	rtrn, err := p.fdb.Transact(func(tr fdb.Transaction) (result interface{}, err error) {
		trOpt := tr.Options()
		trOpt.SetRetryLimit(10)
		titles, err := directory.CreateOrOpen(tr, []string{`titleFiz`}, nil)
		if err != nil {
			return
		}
		list := titles.Sub(`list`)
		hkey1 := titles.Sub(`hkey1`)
		hkey2 := titles.Sub(`hkey2`)
		hkey3 := titles.Sub(`hkey3`)
		unit := titles.Sub(`unit`)
		partner := titles.Sub(`partner`)
		file := titles.Sub(`file`)

		tr.Set(list.Pack(tuple.Tuple{t.Id}), body)

		if h1 := GetHkey1(t); h1 != nil {
			tr.Set(hkey1.Pack(tuple.Tuple{*h1, t.Id}), nil)
		}

		if h2 := GetHkey2(t); h2 != nil {
			tr.Set(hkey2.Pack(tuple.Tuple{*h2, t.Id}), nil)
		}

		if h3 := GetHkey3(t); h3 != nil {
			tr.Set(hkey3.Pack(tuple.Tuple{*h3, t.Id}), nil)
		}

		tr.Set(unit.Pack(tuple.Tuple{t.Unit, t.Id}), nil)
		tr.Set(partner.Pack(tuple.Tuple{t.PartnerId, t.Id}), nil)
		tr.Set(file.Pack(tuple.Tuple{t.FileId, t.Id}), nil)

		return body, err
	})
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(rtrn.([]byte)))
	}
}
