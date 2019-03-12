package title

import (
	"errors"
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/golang/protobuf/proto"
	"github.com/rs/xid"
	"gitlab.xxx.local/xx/components/fdb-proto.git/go"
)

func (p *Class) SetTitle() (msg interface{}, err error) {

	var t = new(fdbProto.DocTitle)
	err = proto.Unmarshal(p.msg.Data, t)
	if err != nil {
		return
	}
	if t.Id == `` {
		t.Id = xid.New().String()
	}

	body, err := proto.Marshal(t)
	if err != nil {
		return
	}

	normalize(t)
	h1 := getHkey1(t)
	h2 := getHkey2(t)
	h3 := getHkey3(t)

	msg, err = p.fdb.Transact(func(tr fdb.Transaction) (i interface{}, e error) {

		if p.tryCount < p.tryCountMax {
			p.tryCount = p.tryCount + 1
		} else {
			tr.Cancel()
			e = errors.New(fmt.Sprintf(`cycle error %v`, p.tryCount))
			return
		}

		titles, e := directory.CreateOrOpen(tr, []string{`titleFiz`}, nil)
		if e != nil {
			return
		}

		list := titles.Sub(`list`)

		if tr.Get(list.Pack(tuple.Tuple{t.Id})).MustGet() != nil {
			e = errors.New(fmt.Sprintf(`Doc with key (%v) already exists`, tuple.Tuple{t.Id}))
			return
		}
		hkey1 := titles.Sub(`hkey1`)
		hkey2 := titles.Sub(`hkey2`)
		hkey3 := titles.Sub(`hkey3`)
		unit := titles.Sub(`unit`)
		partner := titles.Sub(`partner`)
		file := titles.Sub(`file`)

		tr.Set(list.Pack(tuple.Tuple{t.Id}), body)
		if h1 != nil {
			tr.Set(hkey1.Pack(tuple.Tuple{*h1, t.Id}), nil)
		}
		if h2 != nil {
			tr.Set(hkey2.Pack(tuple.Tuple{*h2, t.Id}), nil)
		}
		if h3 != nil {
			tr.Set(hkey3.Pack(tuple.Tuple{*h3, t.Id}), nil)
		}
		tr.Set(unit.Pack(tuple.Tuple{t.Unit, t.Id}), nil)
		tr.Set(partner.Pack(tuple.Tuple{t.PartnerId, t.Id}), nil)
		tr.Set(file.Pack(tuple.Tuple{t.FileId, t.Id}), nil)
		var k = new(fdbProto.Key)
		k.Key = t.Id
		return proto.Marshal(k)
	})
	return
}
