package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Himanshu54/pdf"
)

const dateRex = `(\d{2}-[A-Za-z]{3}-\d{4})`
const amountRex = `([(-]*\d[\d,.]+)\)*`

func readAmount(s string) float64 {
	v, _ := strconv.ParseFloat(strings.ReplaceAll(s, ",", ""), 32)
	return v
}

type Folio struct {
	FolioNo string
	AMC     string
	PAN     string
	KYC     string
	Schemes []Scheme
}
type Scheme struct {
	Scheme       string
	Registrar    string
	Advisor      string
	Scheme_type  string
	Rta_code     string
	Open         float64
	Close        float64
	Value        float64
	Nav          float64
	Valuation    Valuation
	Transactions []Transaction
	Charges      float64
	PL           float64
}
type Transaction struct {
	Date    string
	Unit    float64
	Balance float64
	Amount  float64
	Type    string
	Price   float64
}

type Valuation struct {
	Date  string
	Value float64
	Nav   float64
}
type Info struct {
	Folios []Folio
}

func main() {
	pdf.DebugOn = true
	content, err := readPdf("test.pdf") // Read local pdf file
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(content)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func pw() string {
	fmt.Println("Enter Password: ")
	var ps string
	fmt.Scanln(&ps)
	return ps
}

func readPdf(path string) (Info, error) {
	info := Info{}
	f, err := os.Open(path)
	if err != nil {
		f.Close()
		return info, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return info, err
	}
	r, err := pdf.NewReaderEncrypted(f, fi.Size(), pw)
	if err != nil {
		return info, err
	}
	// remember close file
	defer f.Close()
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return info, err
	}
	buf.ReadFrom(b)
	info.Folios = getFolio(buf.String())
	// fmt.Println(FolioNo_nos)
	return info, err
}
func schemesInfo(folio string, s string) []Scheme {
	ret := []Scheme{}
	r, _ := regexp.Compile(fmt.Sprintf(`PAN:\s(?:OK|NOT\sOK)\s+([0-9A-Z\s]+)\-\s([\(\)A-Z0-9\s/a-z'\-&]+)\(\s+Advisor\s+:\s+([A-Z\-0-9\s]+)\)\sRegistrar\s+\:\s*(CAMS|KFINTECH)\s+Folio\sNo:\s*%s\s*Opening\s*Unit\s*Balance\s*:\s*([\d\.]+)\s*(.*?)NAV\s*on\s*%s:\s*INR\s*%s\s*Valuation\s*on\s*%s:\s*INR\s*%s.*?Closing\s*Unit\sBalance\s*:\s*%s`, folio, dateRex, amountRex, dateRex, amountRex, amountRex))
	// fmt.Printf(`PAN:\s(?:OK|NOT\sOK)\s+([0-9A-Z\s]+)\-\s([\(\)A-Z0-9\s/a-z'\-&]+)\(\s+Advisor\s+:\s+([A-Z\-0-9\s]+)\)\sRegistrar\s+\:\s*(CAMS|KFINTECH)\s+Folio\sNo:\s*%s\s*Opening\s*Unit\s*Balance\s*:\s*([\d\.]+)\s*(.*?)NAV\s*on\s*%s:\s*INR\s*%s\s*Valuation\s*on\s*%s:\s*INR\s*%s.*?Closing\s*Unit\sBalance\s*:\s*%s`, folio, dateRex, amountRex, dateRex, amountRex, amountRex)
	for _, sc := range r.FindAllStringSubmatch(s, -1) {
		p, v, t := transactionInfo(sc[6])
		sch := Scheme{
			Rta_code:  sc[1],
			Scheme:    sc[2],
			Advisor:   sc[3],
			Registrar: sc[4],
			Open:      readAmount(sc[5]),
			Close:     readAmount(sc[11]),
			Nav:       readAmount(sc[8]),
			Valuation: Valuation{
				Date:  sc[9],
				Nav:   readAmount(sc[8]),
				Value: readAmount(sc[10]),
			},
			Value:        v * readAmount(sc[8]),
			Transactions: t,
			Charges:      (v * readAmount(sc[8])) - readAmount(sc[10]),
			PL:           p - readAmount(sc[10]),
		}
		ret = append(ret, sch)
	}
	//  fmt.Println(ret)
	return ret

}

func transactionInfo(s string) (float64, float64, []Transaction) {
	// fmt.Println(s)
	value := 1.0
	price := 0.0
	ret := []Transaction{}
	r, _ := regexp.Compile(fmt.Sprintf(`%s\s*\(?%s\)?\s*%s\s*\(?%s\)?\s*(.*?)\s%s`, dateRex, amountRex, amountRex, amountRex, amountRex))
	// fmt.Printf(`%s\s*\(?%s\)?\s*%s\s*\(?%s\)?\s*(.*)\s(%s)\n`, dateRex, amountRex, amountRex, amountRex, amountRex)
	for _, t := range r.FindAllStringSubmatch(s, -1) {

		op, tp := getTransactionType(t[5])
		value += op * readAmount(t[4])
		price += readAmount(t[2])

		ret = append(ret, Transaction{
			Date:    t[1],
			Amount:  readAmount(t[2]),
			Price:   readAmount(t[3]),
			Unit:    readAmount(t[4]),
			Type:    tp,
			Balance: readAmount(t[6]),
		})
	}
	return price, value, ret
}
func getTransactionType(s string) (float64, string) {

	if strings.Contains(s, "Purchase") {
		return 1, "PURCHASE"
	}
	if strings.Contains(s, "Redemption") {
		return -1, "REDEMPTION"
	}
	if strings.Contains(s, "Switch Over In") {
		return 1, "SWITCH_IN"
	}
	if strings.Contains(s, "Switch Over Out") {
		return -1, "SWITCH_OUT"
	}
	fmt.Println(s)
	return 1, "UNKNOWN"
}
func getFolio(s string) []Folio {
	// fmt.Print(s)
	ret := []Folio{}
	set := make(map[string]Folio)
	r, _ := regexp.Compile(`Balance:?\s*[0-9\.,]*\s*([A-Za-z][A-Za-z\s]+)\s+(?:PAN:\s*\w{5}\d{4}\w)?\s*KYC:\s(OK|NOT\sOK)\s+PAN:\s(OK|NOT\sOK)\s+.*?(?:CAMS|KFINTECH)\s+Folio\sNo:\s([\d//\s]+)\s+`)
	// fmt.Println(r.FindAllStringSubmatch(s, -1)[1:])
	for _, m := range r.FindAllStringSubmatch(s, -1) {
		if _, ok := set[m[1]]; ok {
			continue
		}
		set[m[1]] = Folio{
			AMC:     m[1],
			KYC:     m[2],
			PAN:     m[3],
			FolioNo: m[4],
			Schemes: schemesInfo(m[4], s),
		}
		// fmt.Println(m)
		// break

	}
	for key := range set {
		ret = append(ret, set[key])
	}
	// fmt.Println(ret)
	return ret
}
