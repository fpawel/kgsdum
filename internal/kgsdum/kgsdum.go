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
	products := data.GetProductsOfYearMonth(year, month)

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
