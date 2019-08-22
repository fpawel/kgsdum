package kgsdum

import (
	"fmt"
	"github.com/fpawel/kgsdum/internal/data"
	"strconv"
	"strings"
)

type Table = []Row

type Row = []Cell

type Cell = string

type ValueType int

const (
	vtNone ValueType = iota
	vtOk
	vtError
)

func ProductsOfYearMonthDataTable(year, month int) Table {
	products := GetProductsOfYearMonth(year, month)

	titles := []string{
		"Номер", "День", "Загрузка", "Зав.номер", "Адрес", "U.раб ПГС3",
		"U.раб +20⁰С", "U.cр +20⁰С",
		"U.раб -5⁰С", "U.cр -5⁰С",
		"U.раб +50⁰С", "U.cр +50⁰С",
		"Конц.ПГС1 +20⁰С", "Конц.ПГС4 +20⁰С",
		"Конц.ПГС1 0⁰С", "Конц.ПГС4 0⁰С",
		"Конц.ПГС1 +50⁰С", "Конц.ПГС4 +50⁰С",
		"Конц.ПГС1 +20⁰С", "Конц.ПГС4 +20⁰С",
	}
	r1 := append([]Row{titles}, make([]Row, len(products))...)

	hasCols := make(map[int]struct{})

	for i, p := range products {
		r1[i+1] = make(Row, len(titles))
		r1[i+1][0] = fmt.Sprintf("%d", p.ProductID)
		r1[i+1][1] = fmt.Sprintf("%d", p.Day)
		r1[i+1][2] = fmt.Sprintf("%d", p.PartyID)
		r1[i+1][3] = fmt.Sprintf("%d", p.Serial)
		r1[i+1][4] = fmt.Sprintf("%d", p.Addr)

		r1[i+1][5] = fmtF(p.WorkGas3, 3)
		r1[i+1][6] = fmtF(p.WorkPlus20, 3)
		r1[i+1][7] = fmtF(p.RefPlus20, 3)
		r1[i+1][8] = fmtF(p.WorkMinus5, 3)
		r1[i+1][9] = fmtF(p.RefMinus5, 3)
		r1[i+1][10] = fmtF(p.WorkPlus50, 3)
		r1[i+1][11] = fmtF(p.RefPlus50, 3)
		r1[i+1][12] = fmtF(p.C1Plus20, 3)
		r1[i+1][13] = fmtF(p.C4Plus20, 3)
		r1[i+1][14] = fmtF(p.C1Zero, 3)
		r1[i+1][15] = fmtF(p.C4Zero, 3)
		r1[i+1][16] = fmtF(p.C1Plus50, 3)
		r1[i+1][17] = fmtF(p.C4Plus50, 3)
		r1[i+1][18] = fmtF(p.C1Plus20Ret, 3)
		r1[i+1][19] = fmtF(p.C4Plus20Ret, 3)

		for j := range titles {
			if len(r1[i+1][j]) > 0 {
				hasCols[j] = struct{}{}
			}
		}
	}

	var r2 Table
	for _, row1 := range r1 {
		var row2 Row
		for col := range r1[0] {
			if _, f := hasCols[col]; f {
				row2 = append(row2, row1[col])
			}
		}
		r2 = append(r2, row2)
	}
	return r2
}

type Product struct {
	ProductID   int64    `db:"product_id"`
	PartyID     int64    `db:"party_id"`
	Day         int      `db:"day"`
	Serial      int      `db:"serial_number"`
	Addr        int      `db:"addr"`
	WorkGas3    *float64 `db:"work_gas3"`
	WorkPlus20  *float64 `db:"work_plus20"`
	WorkMinus5  *float64 `db:"work_minus5"`
	WorkPlus50  *float64 `db:"work_plus50"`
	RefPlus20   *float64 `db:"ref_plus20"`
	RefMinus5   *float64 `db:"ref_minus5"`
	RefPlus50   *float64 `db:"ref_plus50"`
	C1Plus20    *float64 `db:"c1_plus20"`
	C1Zero      *float64 `db:"c1_zero"`
	C1Plus50    *float64 `db:"c1_plus50"`
	C1Plus20Ret *float64 `db:"c1_plus20ret"`
	C4Plus20    *float64 `db:"c4_plus20"`
	C4Zero      *float64 `db:"c4_zero"`
	C4Plus50    *float64 `db:"c4_plus50"`
	C4Plus20Ret *float64 `db:"c4_plus20ret"`
}

func GetProductsOfYearMonth(year, month int) (products []Product) {
	err := data.DB.Select(&products, `
SELECT product_id, product.party_id, addr, serial_number,
       work_plus20, ref_plus20,
       work_gas3,
       work_minus5, ref_minus5,
       work_plus50, ref_plus50,
       c1_plus20, c4_plus20,
       c1_zero, c4_zero,
       c1_plus50, c4_plus50,
       c1_plus20ret, c4_plus20ret,
       cast(strftime('%d', party.created_at) AS INTEGER) AS day
FROM product
         INNER JOIN party on product.party_id = party.party_id
WHERE cast(strftime('%Y', party.created_at) AS INTEGER)= ?
  AND cast(strftime('%m', party.created_at) AS INTEGER) = ?
  ORDER BY product.created_at`, year, month)
	if err != nil {
		panic(err)
	}

	return
}

func fmtF(v *float64, prec int) string {
	if v != nil {
		return formatFloat(*v, prec)
	}
	return ""
}

func formatFloat(v float64, prec int) string {
	s := strconv.FormatFloat(v, 'f', prec, 64)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}
