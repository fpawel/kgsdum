package main

import (
	"encoding/base64"
	"flag"
	"github.com/fpawel/gohelp/must"
	"github.com/fpawel/kgsdum/internal/kgsdum"
	"github.com/fpawel/kgsdum/internal/pdf"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatalln("expected 'products.table' or 'pdf' command")
	}

	switch os.Args[1] {

	case "pdf":
		cmd := flag.NewFlagSet("pdf", flag.ExitOnError)
		strProducts := cmd.String("products", "", "coma separated integer identifiers of products to include in pdf, f.e. 1231,456,789")
		if err := cmd.Parse(os.Args[2:]); err != nil {
			cmd.Usage()
			os.Exit(1)
		}
		var products []int64
		for i, s := range strings.Split(*strProducts, ",") {
			id, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				log.Fatalln("products:", "position:", i, ":", err)
			}
			products = append(products, id)
		}
		if len(products) == 0 {
			log.Fatal("products must be set")
		}
		pdf.RunProductIDs(products)

	case "products.table":

		cmd := flag.NewFlagSet("products.table", flag.ExitOnError)
		year := cmd.Int("y", -1, "filter year")
		month := cmd.Int("m", -1, "filter month")
		format := cmd.String("f", "text", "format output (text|base64)")

		if err := cmd.Parse(os.Args[2:]); err != nil {
			cmd.Usage()
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
			cmd.Usage()
			os.Exit(1)
		}
	default:
		log.Fatalln("expected 'products.table' or 'pdf' command")
	}
}
