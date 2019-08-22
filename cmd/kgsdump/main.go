package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/fpawel/gohelp/must"
	"github.com/fpawel/kgsdum/internal/kgsdum"
	"os"
)

func main() {
	data1Cmd := flag.NewFlagSet("d1", flag.ExitOnError)
	year := data1Cmd.Int("y", 2019, "filter year")
	month := data1Cmd.Int("m", 1, "filter month")
	format := data1Cmd.String("f", "text", "format output (text|base64)")

	if len(os.Args) < 2 {
		fmt.Println("expected 'd1' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "d1":
		if err := data1Cmd.Parse(os.Args[2:]); err != nil {
			flag.Usage()
			panic(err)
		}
		bs := must.MarshalJSON(kgsdum.ProductsOfYearMonthDataTable(*year, *month))
		switch *format {
		case "text":
			must.Write(os.Stdout, bs)
		case "base64":
			must.Write(base64.NewEncoder(base64.StdEncoding, os.Stdout), bs)
		default:
			panic("wrong format: " + *format)

		}

	default:
		fmt.Println("expected 'd1' subcommand")
		os.Exit(1)
	}
}
