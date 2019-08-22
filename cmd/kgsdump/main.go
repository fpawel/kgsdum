package main

import (
	"encoding/base64"
	"flag"
	"github.com/fpawel/gohelp/must"
	"github.com/fpawel/kgsdum/internal/kgsdum"
	"log"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatalln("expected 'products.table' command")
	}

	switch os.Args[1] {

	case "products.table":

		productsTableCmd := flag.NewFlagSet("products.table", flag.ExitOnError)
		year := productsTableCmd.Int("y", -1, "filter year")
		month := productsTableCmd.Int("m", -1, "filter month")
		format := productsTableCmd.String("f", "text", "format output (text|base64)")

		if err := productsTableCmd.Parse(os.Args[2:]); err != nil {
			productsTableCmd.Usage()
			os.Exit(1)
		}

		if *year == -1 {
			log.Fatalln("filter year must be set")
		}
		if *month == -1 {
			log.Fatalln("filter month must be set")
		}

		bs := must.MarshalJSON(kgsdum.ProductsOfYearMonthDataTable(*year, *month))
		switch *format {
		case "text":
			must.Write(os.Stdout, bs)
		case "base64":
			if _, err := os.Stdout.WriteString(base64.StdEncoding.EncodeToString(bs)); err != nil {
				panic(err)
			}
		default:
			log.Println("wrong format:", *format)
			productsTableCmd.Usage()
			os.Exit(1)
		}

	default:
		log.Fatalln("expected 'products.table' command")
	}
}
