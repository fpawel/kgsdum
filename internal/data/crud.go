package data

import "strings"

func GetProductsOfYearMonth(year, month int) (products []Product) {
	err := DB.Select(&products, `
SELECT product_id, product.party_id, addr, serial_number,
       cast(strftime('%Y', party.created_at) AS INTEGER) AS year,
       cast(strftime('%m', party.created_at) AS INTEGER) AS month,
       cast(strftime('%d', party.created_at) AS INTEGER) AS day,
       
       work_plus20, ref_plus20,
       work_gas3,
       work_minus5, ref_minus5,
       work_plus50, ref_plus50,
       c1_plus20, c4_plus20,
       c1_zero, c4_zero,
       c1_plus50, c4_plus50,
       c1_plus20ret, c4_plus20ret,       
       100 * (c1_plus20 - pgs1) / (0.1 + pgs1 * 0.12) AS err1_plus20,
       100 * (c4_plus20 - pgs1) / (0.1 + pgs4 * 0.12) AS err4_plus20,
       100 * (c1_zero - pgs1) / (0.1 + pgs1 * 0.12) AS err1_zero,
       100 * (c4_zero - pgs1) / (0.1 + pgs4 * 0.12) AS err4_zero,
       100 * (c1_plus50 - pgs1) / (0.1 + pgs1 * 0.12) AS err1_plus50,
       100 * (c4_plus50 - pgs1) / (0.1 + pgs4 * 0.12) AS err4_plus50,
       100 * (c1_plus20ret - pgs1) / (0.1 + pgs1 * 0.12) AS err1_plus20ret,
       100 * (c4_plus20ret - pgs1) / (0.1 + pgs4 * 0.12) AS err4_plus20ret,
       pgs1,
       pgs2,
       pgs3,
       pgs4
FROM product
         INNER JOIN party on product.party_id = party.party_id
WHERE year = ?
  AND month = ?
  ORDER BY product.created_at`, year, month)
	if err != nil {
		panic(err)
	}

	return
}

func GetProductsByID(productIDs []int64) (products []Product) {
	args := make([]interface{}, len(productIDs))
	for i := range args {
		args[i] = productIDs[i]
	}
	err := DB.Select(&products, `
SELECT product_id, product.party_id, addr, serial_number,
       
       cast(strftime('%Y', party.created_at) AS INTEGER) AS year,
       cast(strftime('%m', party.created_at) AS INTEGER) AS month,
       cast(strftime('%d', party.created_at) AS INTEGER) AS day,
       
       work_plus20, ref_plus20,
       work_gas3,
       work_minus5, ref_minus5,
       work_plus50, ref_plus50,
       c1_plus20, c4_plus20,
       c1_zero, c4_zero,
       c1_plus50, c4_plus50,
       c1_plus20ret, c4_plus20ret,
       100 * (c1_plus20 - pgs1) / (0.1 + pgs1 * 0.12) AS err1_plus20,
       100 * (c4_plus20 - pgs1) / (0.1 + pgs4 * 0.12) AS err4_plus20,
       100 * (c1_zero - pgs1) / (0.1 + pgs1 * 0.12) AS err1_zero,
       100 * (c4_zero - pgs1) / (0.1 + pgs4 * 0.12) AS err4_zero,
       100 * (c1_plus50 - pgs1) / (0.1 + pgs1 * 0.12) AS err1_plus50,
       100 * (c4_plus50 - pgs1) / (0.1 + pgs4 * 0.12) AS err4_plus50,
       100 * (c1_plus20ret - pgs1) / (0.1 + pgs1 * 0.12) AS err1_plus20ret,
       100 * (c4_plus20ret - pgs1) / (0.1 + pgs4 * 0.12) AS err4_plus20ret,
       pgs1,
       pgs2,
       pgs3,
       pgs4
FROM product
         INNER JOIN party on product.party_id = party.party_id`+
		" WHERE product_id IN (?"+strings.Repeat(",?", len(productIDs)-1)+")", args...)
	if err != nil {
		panic(err)
	}
	return
}
