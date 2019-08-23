package data

import (
	"fmt"
	"testing"
)

func TestGetProductsByID(t *testing.T) {
	fmt.Println(GetProductsByID([]int64{1, 5, 6}))
}
