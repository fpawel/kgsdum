package pdf

import (
	"fmt"
	"github.com/fpawel/gohelp/must"
	"github.com/fpawel/gohelp/winapp"
	"github.com/fpawel/kgsdum/internal/data"
	"github.com/jung-kurt/gofpdf"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func RunProductIDs(productIDs []int64) {

	d := gofpdf.New("P", "mm", "A4",
		filepath.Join(filepath.Dir(os.Args[0]), "fonts"))
	d.AddFont("RobotoCondensed", "", "RobotoCondensed-Regular.json")
	d.AddFont("RobotoCondensed", "B", "RobotoCondensed-Bold.json")
	d.AddFont("RobotoCondensed", "I", "RobotoCondensed-Italic.json")
	d.AddFont("RobotoCondensed", "BI", "RobotoCondensed-BoldItalic.json")
	d.AddFont("RobotoCondensed-Light", "", "RobotoCondensed-Light.json")
	d.AddFont("RobotoCondensed-Light", "I", "RobotoCondensed-LightItalic.json")
	d.UnicodeTranslatorFromDescriptor("cp1251")
	d.SetLineWidth(.1)
	d.SetFillColor(225, 225, 225)
	d.SetDrawColor(169, 169, 169)

	d.AddPage()

	products := data.GetProductsByID(productIDs)
	for i, p := range products {
		pdfProduct(d, p)
		if i == len(products)-1 {
			break
		}
		if (i+1)%4 == 0 {
			d.AddPage()
			continue
		}
		d.Ln(5)
		w := d.GetLineWidth()

		d.SetLineWidth(0.25)
		d.SetDashPattern([]float64{2}, 2)

		d.MoveTo(d.GetX(), d.GetY())
		d.LineTo(d.GetX()+120, d.GetY())
		d.DrawPath("D")

		d.SetLineWidth(w)
		d.SetDashPattern([]float64{}, 0)

		d.Ln(5)
	}
	saveAndShowDoc(d)
}

func pdfProduct(d *gofpdf.Fpdf, p data.Product) {
	tr := d.UnicodeTranslatorFromDescriptor("cp1251")
	sentence2 := func(familyStr, fontStyleStr string, fontSize float64, h float64, s string) {
		d.SetFont(familyStr, fontStyleStr, fontSize)
		d.CellFormat(d.GetStringWidth(tr(s)), h, tr(s), "", 0, "", false, 0, "")
	}

	sentence := func(fontStyleStr, s string) {
		familyStr := "RobotoCondensed-Light"
		if fontStyleStr == "B" {
			familyStr = "RobotoCondensed"
		}
		sentence2(familyStr, fontStyleStr, fontSize1, lineSpace1, s)
	}

	sentencef := func(fontStyleStr, format string, args ...interface{}) {
		sentence(fontStyleStr, fmt.Sprintf(format, args...))
	}

	//pageWidth, _ := d.GetPageSize()

	//width := pageWidth - spaceX*2

	d.SetX(spaceX)
	d.SetFont("RobotoCondensed", "B", 11)
	d.CellFormat(120, 8, tr("Паспорт блока оптического ИБЯЛ.418414.086"), "", 1, "C", false, 0, "")
	d.SetX(spaceX)

	sentence("", "Дата изготовления: ")
	sentencef("B", "%02d.%02d.%d ", p.Day, p.Month, p.Year)
	d.SetX(d.GetX() + 5)

	sentence("", "Заводской номер: ")
	sentencef("B", "%d ", p.Serial)

	d.Ln(-1)

	sentence("", "Партия: ")
	sentencef("B", "%d", p.PartyID)
	d.SetX(d.GetX() + 5)

	sentence("", "ID: ")
	sentencef("B", "%d ", p.ProductID)
	d.SetX(d.GetX() + 5)

	sentence("", "Полезный сигнал: ")
	if p.C1Plus20 != nil && p.WorkGas3 != nil && *p.C1Plus20 != 0 {
		sentencef("B", "%s", formatFloat(*p.C1Plus20 - *p.WorkGas3 / *p.C1Plus20, 3))
	} else {
		sentence("B", "???")
	}
	d.Ln(lineSpace1)

	type Cell struct {
		Align, Style, Text string
	}

	colWidths := []float64{40, 12, 12, 12}

	for _, row := range [][]Cell{
		{
			Cell{"L", "I", "Температура,\"С"},
			Cell{"C", "", "-5\"С"},
			Cell{"C", "", "+20\"С"},
			Cell{"C", "", "+50\"С"},
		},
		{
			Cell{"L", "I", "Сигнал рабочего канала"},
			Cell{"R", "B", formatNullFloat64(p.WorkMinus5, 3)},
			Cell{"R", "B", formatNullFloat64(p.WorkPlus20, 3)},
			Cell{"R", "B", formatNullFloat64(p.WorkPlus50, 3)},
		},
		{
			Cell{"L", "I", "Сигнал сравнительного канала"},
			Cell{"R", "B", formatNullFloat64(p.RefMinus5, 3)},
			Cell{"R", "B", formatNullFloat64(p.RefPlus20, 3)},
			Cell{"R", "B", formatNullFloat64(p.RefPlus50, 3)},
		},
	} {
		d.SetX(spaceX)
		for col, cell := range row {
			if cell.Style == "B" {
				d.SetFont("RobotoCondensed", cell.Style, 8)
			} else {
				d.SetFont("RobotoCondensed-Light", cell.Style, 8)
			}
			d.CellFormat(colWidths[col], 3.5, tr(cell.Text), "1", 0, cell.Align, false, 0, "")
		}
		d.Ln(-1)
	}

	d.Ln(-1)

	colWidths = []float64{30, 35, 12, 12, 12, 12}

	for _, row := range [][]Cell{
		{
			Cell{"L", "I", "ГСО-ПГС, % об.дол."},
			Cell{"L", "I", "Погрешность, % об.дол."},
			Cell{"C", "I", "+20\"С"},
			Cell{"C", "", "0\"С"},
			Cell{"C", "", "+50\"С"},
			Cell{"C", "", "+20\"С"},
		},
		{
			Cell{"L", "", "ПГС1, азот"},
			Cell{"L", "", "0.1"},
			Cell{"R", "B", formatNullFloat64(p.C1Plus20, 3)},
			Cell{"R", "B", formatNullFloat64(p.C1Zero, 3)},
			Cell{"R", "B", formatNullFloat64(p.C1Plus50, 3)},
			Cell{"R", "B", formatNullFloat64(p.C1Plus20Ret, 3)},
		},
		{
			Cell{"L", "", fmt.Sprintf("ПГС4 = %v", p.Pgs4)},
			Cell{"L", "", fmt.Sprintf("0.1 + %v * 0.12 = %v", p.Pgs4, 0.1+p.Pgs4*0.12)},
			Cell{"R", "B", formatNullFloat64(p.C4Plus20, 3)},
			Cell{"R", "B", formatNullFloat64(p.C4Zero, 3)},
			Cell{"R", "B", formatNullFloat64(p.C4Plus50, 3)},
			Cell{"R", "B", formatNullFloat64(p.C4Plus20Ret, 3)},
		},
	} {
		d.SetX(spaceX)
		for col, cell := range row {
			if cell.Style == "B" {
				d.SetFont("RobotoCondensed", cell.Style, 8)
			} else {
				d.SetFont("RobotoCondensed-Light", cell.Style, 8)
			}
			d.CellFormat(colWidths[col], 3.5, tr(cell.Text), "1", 0, cell.Align, false, 0, "")
		}
		d.Ln(-1)
	}
	d.Ln(-1)
	sentence("", "Подпись:")
	d.SetX(d.GetX() + 5)
	d.CellFormat(40, lineSpace1, "", "B", 1, "", false, 0, "")
}

func saveAndShowDoc(d *gofpdf.Fpdf) {

	dir := filepath.Join(filepath.Dir(os.Args[0]), "pdf")
	if err := winapp.EnsuredDirectory(dir); err != nil {
		panic(err)
	}
	_ = clearDir(dir)

	tmpFile, err := ioutil.TempFile(dir, "")
	if err != nil {
		panic(err)
	}
	must.Close(tmpFile)
	must.Remove(tmpFile.Name())

	pdfFileName := tmpFile.Name() + ".pdf"

	if err := d.OutputFileAndClose(pdfFileName); err != nil {
		panic(err)
	}
	if err := exec.Command("explorer.exe", pdfFileName).Start(); err != nil {
		panic(err)
	}
}

func clearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func formatNullFloat64(v *float64, precision int) string {
	if v != nil {
		return formatFloat(*v, precision)
	}
	return ""
}

func formatFloat(v float64, precision int) string {
	s := strconv.FormatFloat(v, 'f', precision, 64)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

const (
	fontSize1  = 9
	lineSpace1 = 6
	spaceX     = 10.
)
