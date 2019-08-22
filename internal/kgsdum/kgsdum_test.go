package kgsdum

import (
	"fmt"
	"testing"
)

func TestProductsDataTable(t *testing.T) {

	for _, p := range ProductsOfYearMonthDataTable(2019, 4) {
		fmt.Printf("%+v\n", p)
	}

}
