package kgsdum

import (
	"fmt"
	"github.com/fpawel/kgsdum/internal/data"
	"math/rand"
	"testing"
	"time"
)

func TestProductsDataTable(t *testing.T) {

	for _, p := range ProductsOfYearMonthDataTable(2019, 4) {
		fmt.Printf("%+v\n", p)
	}

}

func TestCreateDB(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for y := 2018; y <= 2019; y++ {
		for m := time.Month(3); m <= 5; m++ {
			for d := 1; d <= 3; d++ {
				fmt.Printf("%02d:%02d:%02d\n", d, m, y)
				t := time.Date(y, m, d, 16, 12, 0, 0, time.Now().Location())
				r, err := data.DB.Exec(`INSERT INTO party (created_at) VALUES (?)`, t)
				if err != nil {
					panic(err)
				}
				partyID, err := r.LastInsertId()
				if err != nil {
					panic(err)
				}
				f := func() float64 {
					return rand.Float64() * 100
				}
				for i := 0; i < 20; i++ {
					t := t.Add(time.Second * time.Duration(i+1))
					if d == 1 {
						data.DB.MustExec(`
INSERT INTO product( created_at, party_id, serial_number, addr, production, work_plus20, ref_plus20, work_plus50, ref_plus50, c1_plus20, c4_plus20, c1_zero, c4_zero, c1_plus50, c4_plus50, c1_plus20ret, c4_plus20ret) 
VALUES (?, ?,?,?,TRUE, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,? )`,
							t, partyID, i+100, i+1, f(), f(), f(), f(), f(), f(), f(), f(), f(), f(), f(), f(), f(), f(), f())
					} else {
						data.DB.MustExec(`
INSERT INTO product(  created_at, party_id, serial_number, addr, production, work_plus20, ref_plus20, work_gas3, work_minus5, ref_minus5, work_plus50, ref_plus50, c1_plus20, c4_plus20, c1_zero, c4_zero ) 
VALUES (?, ?,?,?,TRUE, ?, ?, ?, ?, ?, ?, ?,  ?, ?, ?,? )`,
							t, partyID, i+100, i+1, f(), f(), f(), f(), f(), f(), f(), f(), f(), f(), f(), f(), f(), f(), f())
					}

				}
			}
		}
	}
}
