package program

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/rs/xid"
	"strconv"
	"time"
)

func (p *Program) GenerateTitle() Title {
	Birthday := generateRandomDate("1916-08-01", "2017-08-22")
	DocNo := strconv.Itoa(randomdata.Number(1000000000, 9999999999))

	sex := map[int]int{
		0: 1,
		1: 2,
		2: 9,
	}

	return Title{
		Id:         xid.New().String(),
		FileId:     int64(randomdata.Number(0, 9999999999)),
		PartnerId:  int64(randomdata.Number(0, 9999999999)),
		Unit:       int64(randomdata.Number(0, 9999999999)),
		LastName:   randomdata.LastName(),
		FirstName:  randomdata.FirstName(randomdata.RandomGender),
		MiddleName: randomdata.Title(randomdata.RandomGender),
		Birthday:   Birthday,
		DocNo:      DocNo,
		DocType:    randomdata.Number(1, 16),
		Sex:        sex[randomdata.Number(0, 2)],
		Address: addrs{
			Reg: &addr{
				ZipCode:   randomdata.PostalCode("RU"),
				Country:   randomdata.Country(0),
				Region:    "",
				City:      randomdata.City(),
				District:  randomdata.MacAddress(),
				Statement: "",
				Street:    randomdata.Street(),
				House:     strconv.Itoa(randomdata.Number(99)),
				Block:     strconv.Itoa(randomdata.Number(99)),
				Build:     strconv.Itoa(randomdata.Number(99)),
				Flat:      strconv.Itoa(randomdata.Number(99)),
			}, Fact: &addr{
				ZipCode:   randomdata.PostalCode("RU"),
				Country:   randomdata.Country(0),
				Region:    "",
				City:      randomdata.City(),
				District:  randomdata.MacAddress(),
				Statement: "",
				Street:    randomdata.Street(),
				House:     strconv.Itoa(randomdata.Number(99)),
				Block:     strconv.Itoa(randomdata.Number(99)),
				Build:     strconv.Itoa(randomdata.Number(99)),
				Flat:      strconv.Itoa(randomdata.Number(99)),
			},
		},
	}
}

func generateRandomDate(from string, to string) time.Time {

	t, _ := time.Parse(randomdata.DateOutputLayout, randomdata.FullDateInRange(from, to))

	return t
}
