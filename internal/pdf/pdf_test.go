package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRunProductIDs(t *testing.T) {
	fmt.Println(filepath.Dir(os.Args[0]))
	RunProductIDs([]int64{327, 328, 329, 330, 331, 332, 333, 334, 335, 336, 337})
}
