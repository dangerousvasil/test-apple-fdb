package program

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
)

func toHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))

	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func GetHkey1(t Title) *string {
	var res = struct {
		LastName   string
		FirstName  string
		MiddleName string
		DocType    int
		DocNo      string
	}{
		LastName:   t.LastName,
		FirstName:  t.FirstName,
		MiddleName: t.MiddleName,
		DocType:    t.DocType,
		DocNo:      t.DocNo,
	}
	b, _ := json.Marshal(res)
	s := toHash(string(b))
	return &s
}
func GetHkey2(t Title) *string {
	var res = struct {
		BirthDay string
		DocType  int
		DocNo    string
	}{
		BirthDay: t.Birthday.Format(`2006-02-01`),
		DocType:  t.DocType,
		DocNo:    t.DocNo,
	}
	b, _ := json.Marshal(res)
	s := toHash(string(b))
	return &s
}
func GetHkey3(t Title) *string {
	if t.DocType == 1 {
		var res = struct {
			LastName   string
			FirstName  string
			MiddleName string
			DocReg     string
		}{
			LastName:   t.LastName,
			FirstName:  t.FirstName,
			MiddleName: t.MiddleName,
			DocReg:     string([]rune(t.DocNo)[:2]),
		}
		b, _ := json.Marshal(res)
		s := toHash(string(b))
		return &s
	} else {
		return nil
	}
}
