package data

type Product struct {
	ProductID     int64    `db:"product_id"`
	PartyID       int64    `db:"party_id"`
	Year          int      `db:"year"`
	Month         int      `db:"month"`
	Day           int      `db:"day"`
	Serial        int      `db:"serial_number"`
	Addr          int      `db:"addr"`
	WorkGas3      *float64 `db:"work_gas3"`
	WorkPlus20    *float64 `db:"work_plus20"`
	WorkMinus5    *float64 `db:"work_minus5"`
	WorkPlus50    *float64 `db:"work_plus50"`
	RefPlus20     *float64 `db:"ref_plus20"`
	RefMinus5     *float64 `db:"ref_minus5"`
	RefPlus50     *float64 `db:"ref_plus50"`
	C1Plus20      *float64 `db:"c1_plus20"`
	C1Zero        *float64 `db:"c1_zero"`
	C1Plus50      *float64 `db:"c1_plus50"`
	C1Plus20Ret   *float64 `db:"c1_plus20ret"`
	C4Plus20      *float64 `db:"c4_plus20"`
	C4Zero        *float64 `db:"c4_zero"`
	C4Plus50      *float64 `db:"c4_plus50"`
	C4Plus20Ret   *float64 `db:"c4_plus20ret"`
	Err1Plus20    *float64 `db:"err1_plus20"`
	Err4Plus20    *float64 `db:"err4_plus20"`
	Err1Zero      *float64 `db:"err1_zero"`
	Err4Zero      *float64 `db:"err4_zero"`
	Err1Plus50    *float64 `db:"err1_plus50"`
	Err4Plus50    *float64 `db:"err4_plus50"`
	Err1Plus20Ret *float64 `db:"err1_plus20ret"`
	Err4Plus20Ret *float64 `db:"err4_plus20ret"`
	Pgs1          float64  `db:"pgs1"`
	Pgs2          float64  `db:"pgs2"`
	Pgs3          float64  `db:"pgs3"`
	Pgs4          float64  `db:"pgs4"`
}
